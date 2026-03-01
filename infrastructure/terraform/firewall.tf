resource "google_compute_firewall" "allow_ssh" {
  name    = "prism-allow-ssh"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["${var.allowed_ip}/32"]
  target_tags   = ["prism"]
}
