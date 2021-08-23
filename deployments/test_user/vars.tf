

variable "custom_tags" {
  description = "A map of key value pairs that represents custom tags to apply to taggable resources"
  type        = map(string)
  default     = {}
}
