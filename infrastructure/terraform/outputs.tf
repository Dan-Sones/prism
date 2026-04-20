output "app_internal_ip" {
  value = google_compute_instance.app.network_interface[0].network_ip
}

output "app_external_ip" {
  value = google_compute_instance.app.network_interface[0].access_config[0].nat_ip
}

output "gatling_ip" {
  value = google_compute_instance.gatling.network_interface[0].access_config[0].nat_ip
}


output "gatling_vm_name" {
  value = google_compute_instance.gatling.name
}

