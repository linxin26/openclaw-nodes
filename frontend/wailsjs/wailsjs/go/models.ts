export namespace wails {
	
	export class AboutInfo {
	    deviceId: string;
	    publicKey: string;
	    version: string;
	    platform: string;
	    hostname: string;
	    goVersion: string;
	    arch: string;
	    dataDir: string;
	    protocolVersion: number;
	
	    static createFrom(source: any = {}) {
	        return new AboutInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.deviceId = source["deviceId"];
	        this.publicKey = source["publicKey"];
	        this.version = source["version"];
	        this.platform = source["platform"];
	        this.hostname = source["hostname"];
	        this.goVersion = source["goVersion"];
	        this.arch = source["arch"];
	        this.dataDir = source["dataDir"];
	        this.protocolVersion = source["protocolVersion"];
	    }
	}
	export class ActivityEntry {
	    timestamp: number;
	    event: string;
	    level: string;
	
	    static createFrom(source: any = {}) {
	        return new ActivityEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.event = source["event"];
	        this.level = source["level"];
	    }
	}
	export class CapabilityInfo {
	    key: string;
	    name: string;
	    description: string;
	    enabled: boolean;
	    available: boolean;
	    permission: string;
	    reason?: string;
	    commands: string[];
	    dependencies: string[];
	    healthy: boolean;
	    tier: number;
	
	    static createFrom(source: any = {}) {
	        return new CapabilityInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.enabled = source["enabled"];
	        this.available = source["available"];
	        this.permission = source["permission"];
	        this.reason = source["reason"];
	        this.commands = source["commands"];
	        this.dependencies = source["dependencies"];
	        this.healthy = source["healthy"];
	        this.tier = source["tier"];
	    }
	}
	export class CapabilityOption {
	    provider?: string;
	    path?: string;
	
	    static createFrom(source: any = {}) {
	        return new CapabilityOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.path = source["path"];
	    }
	}
	export class CapabilityState {
	    enabled: boolean;
	    available: boolean;
	    permission: string;
	    reason?: string;
	
	    static createFrom(source: any = {}) {
	        return new CapabilityState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.available = source["available"];
	        this.permission = source["permission"];
	        this.reason = source["reason"];
	    }
	}
	export class Config {
	    gateway: string;
	    port: number;
	    token: string;
	    tls: boolean;
	    discovery: string;
	    capabilities: {[key: string]: boolean};
	    capabilityOptions: {[key: string]: CapabilityOption};
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gateway = source["gateway"];
	        this.port = source["port"];
	        this.token = source["token"];
	        this.tls = source["tls"];
	        this.discovery = source["discovery"];
	        this.capabilities = source["capabilities"];
	        this.capabilityOptions = this.convertValues(source["capabilityOptions"], CapabilityOption, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConnectionStatus {
	    status: string;
	    gateway: string;
	    tls: boolean;
	    uptimeMs: number;
	    retryCount: number;
	    retryDelayMs: number;
	    protocolVersion: number;
	    capabilities: {[key: string]: CapabilityState};
	
	    static createFrom(source: any = {}) {
	        return new ConnectionStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.gateway = source["gateway"];
	        this.tls = source["tls"];
	        this.uptimeMs = source["uptimeMs"];
	        this.retryCount = source["retryCount"];
	        this.retryDelayMs = source["retryDelayMs"];
	        this.protocolVersion = source["protocolVersion"];
	        this.capabilities = this.convertValues(source["capabilities"], CapabilityState, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DeviceInfo {
	    deviceId: string;
	    platform: string;
	    hostname: string;
	    mode: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new DeviceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.deviceId = source["deviceId"];
	        this.platform = source["platform"];
	        this.hostname = source["hostname"];
	        this.mode = source["mode"];
	        this.version = source["version"];
	    }
	}
	export class InvokeResult {
	    success: boolean;
	    data: {[key: string]: any};
	    error: string;
	    durationMs: number;
	
	    static createFrom(source: any = {}) {
	        return new InvokeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	        this.durationMs = source["durationMs"];
	    }
	}
	export class LogEntry {
	    timestamp: number;
	    level: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = source["timestamp"];
	        this.level = source["level"];
	        this.message = source["message"];
	    }
	}
	export class LogFilter {
	    levels: string[];
	    search: string;
	    limit: number;
	    offset: number;
	
	    static createFrom(source: any = {}) {
	        return new LogFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.levels = source["levels"];
	        this.search = source["search"];
	        this.limit = source["limit"];
	        this.offset = source["offset"];
	    }
	}
	export class TestResult {
	    success: boolean;
	    latencyMs: number;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new TestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.latencyMs = source["latencyMs"];
	        this.error = source["error"];
	    }
	}

}

