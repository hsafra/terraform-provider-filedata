// Copyright (c) Harel Safra

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
	"strconv"
	"strings"
	"terraform-provider-filedata/api"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &File{}
var _ resource.ResourceWithImportState = &File{}

func NewFile() resource.Resource {
	return &File{}
}

// File defines the resource implementation.
type File struct {
	basePath string
}

// FileModel describes the resource data model.
type FileModel struct {
	File_name types.String   `tfsdk:"file_name"`
	Lines     []types.String `tfsdk:"lines"`
}

func (r *File) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *File) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "File data resource",

		Attributes: map[string]schema.Attribute{
			"file_name": schema.StringAttribute{
				Description: "File name",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9]+$`),
						"must contain only lowercase alphanumeric characters",
					),
				},
			},
			"lines": schema.ListAttribute{
				Description: "Lines to write to the file",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(2),
				},
			},
		},
	}
}

func (r *File) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	basePath, ok := req.ProviderData.(string)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.basePath = basePath
}

func (r *File) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FileModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fullName := r.basePath + "/" + data.File_name.ValueString()

	// iterate over data.lines and write to file
	for i := 0; i < len(data.Lines); i++ {
		tflog.Trace(ctx, "writing line number "+strconv.Itoa(i)+" to file "+fullName+" with value "+data.Lines[i].ValueString())
		err := api.WriteLine(fullName, i+1, data.Lines[i].ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Write Error", fmt.Sprintf("Unable to write line %d, got error: %s", i, err))
			return
		}
	}

	tflog.Trace(ctx, "file written successfully")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *File) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FileModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fullName := r.basePath + "/" + data.File_name.ValueString()

	lineCount, err := api.LineCount(fullName)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read file, got error: %s", err))
		}
		return
	}

	data.Lines = make([]types.String, 0)
	// read the lines from the file and return them
	for i := 0; i < lineCount; i++ {
		line, err := api.ReadLine(fullName, i+1)
		if err != nil {
			resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read line %d, got error: %s", i, err))
			return
		}
		data.Lines = append(data.Lines, types.StringValue(line))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *File) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, data FileModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.File_name = state.File_name

	fullName := r.basePath + "/" + data.File_name.ValueString()

	// find differences between plan and state and update the file
	planLinecount := len(plan.Lines)
	stateLinecount := len(state.Lines)

	data.Lines = make([]types.String, 0)

	for i := 0; i < planLinecount; i++ {
		data.Lines = append(data.Lines, plan.Lines[i])
		if i >= stateLinecount || plan.Lines[i] != state.Lines[i] {
			tflog.Trace(ctx, "updating line number "+strconv.Itoa(i)+" to file "+fullName+" with value "+plan.Lines[i].ValueString())
			err := api.WriteLine(fullName, i+1, plan.Lines[i].ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Write Error", fmt.Sprintf("Unable to write line %d, got error: %s", i, err))
				return
			}
		}
	}

	if planLinecount < stateLinecount {
		err := api.TrimFile(fullName, planLinecount)
		if err != nil {
			resp.Diagnostics.AddError("Trim Error", fmt.Sprintf("Unable to trim file, got error: %s", err))
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *File) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FileModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	fullName := r.basePath + "/" + data.File_name.ValueString()

	err := api.RemoveFile(fullName)

	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete file, got error: %s", err))
		return
	}
}

func (r *File) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
