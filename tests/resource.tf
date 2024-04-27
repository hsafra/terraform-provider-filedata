resource "filedata_file" "file1" {
  file_name = "file1"
  lines = [
    "line1",
    "line2",
    "line3",
  ]
}


resource "filedata_file" "file2" {
  file_name = "file2"
  lines = [
    "aaaaaaaaaaaa",
    "bbbbbbbb",
    "cccccccccccccccccc",
  ]
}
