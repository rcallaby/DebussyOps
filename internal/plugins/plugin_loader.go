package plugins

// Plugin loader stub. For advanced use you can implement Go plugin loading or a gRPC-based plugin system.
// This file is intentionally left as a template for plugin patterns.

// Suggested approaches:
// - Use Go plugins (.so) for unix builds (plugin package) - limited portability.
// - Use gRPC/HTTP to connect to external agent services (recommended).
// - Use a plugin registry (signed manifests) to discover remote agents.