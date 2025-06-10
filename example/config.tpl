app:
  name: "${app_name}"
  environment: "${environment}"
  port: ${port}
  enable_auth: ${enable_auth}

network:
  allowed_ips:
%{ for ip in allowed_ips ~}
    - ${ip}
%{ endfor ~}
