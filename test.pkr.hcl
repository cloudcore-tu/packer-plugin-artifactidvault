packer {
  required_plugins {
    amazon = {
      version = ">= 1.1.1"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "ubuntu" {
  ami_name      = "packer-plugin-test"
  instance_type = "t2.micro"
  region        = "ap-northeast-1"
  source_ami_filter {
    filters = {
      name = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username          = "ubuntu"
  force_deregister      = true
  force_delete_snapshot = true
  launch_block_device_mappings {
    device_name           = "/dev/sda1"
    volume_size           = 20
    delete_on_termination = true
  }
  tags = {
    Name        = "packer-plugin-test"
    CreateDate  = "2022-05-20"
    Environment = "Test"
    For         = "packer-plugin test"
  }
}

build {
  name = "packer-plugin-test-image-builder"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner shell {
    inline = ["echo hello"]
  }

  post-processor artifactidvault-ssm {
    parameter-name = "/packer-plugin-artifactidvault/test"
    region         = "us-east-1"
    matcher        = "ami-[0-9a-fA-F]+"
    overwrite      = true
  }
}
