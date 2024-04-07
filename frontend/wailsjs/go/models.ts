export namespace main {
	
	export class Model {
	    name: string;
	    size: number;
	    download: boolean;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Model(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.size = source["size"];
	        this.download = source["download"];
	        this.active = source["active"];
	    }
	}

}

