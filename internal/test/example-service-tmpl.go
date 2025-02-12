package internal_test

const (
	ServiceComposeTemplate string = `
services:
  example:
    image: {{.ServiceTag}}
    container_name: {{.ServiceContainerName}} # Container names are provided by Hyperdrive to handle conflicts 
    restart: unless-stopped
    ports:
      - "127.0.0.1:8085:8085/tcp" # Restricted to localhost outside of Docker
    volumes:
      - {{.CfgDir}}:{{.CfgDir}}:ro # The module's config directory
      - {{.LogDir}}:{{.LogDir}} # The module's log directory
    command:
      - "--config-file"
      - "{{.CfgDir}}/service-cfg.yaml" # The module's config directory
      - "--log-dir"
      - "{{.LogDir}}" # The module's log directory
      - --ip
      - "0.0.0.0" # Open to all Docker traffic
      - --port
      - "8085" # The port the module's server will listen on
    networks:
      - net
      - {{.ProjectName}}_net
    cap_drop:
      - all # The container will run as root but this will drop all capabilities it doesn't need
    cap_add:
      - dac_override # Needed to access files owned by the user instead of root
    security_opt:
      - no-new-privileges
networks:
  net:
  {{.ProjectName}}_net:
    external: true
`
)
