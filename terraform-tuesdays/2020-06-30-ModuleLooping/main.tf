variable "mod_count" {
  type    = number
  default = 2
}

module "files" {
  count  = var.mod_count
  source = "./file-creator"

  content  = "Hi! My name is ${count.index}"
  filename = "my_file_${count.index}.txt"
}

variable "file_foreach" {
  default = {
    arthur         = "dent"
    tricia         = "mcmillian"
    zaphod         = "beeblebrox"
    ford           = "prefect"
    slartibartfast = "????"
  }
}

module "files_foreach" {
  source = "./file-creator"

  for_each = var.file_foreach

  content  = each.value
  filename = "${each.key}.txt"

}