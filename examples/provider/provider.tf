terraform {
  required_providers {
    statuspage = {
      source = "winebarrel/statuspage"
    }
  }
}

provider "statuspage" {
  api_key = "your-api-key"
}
