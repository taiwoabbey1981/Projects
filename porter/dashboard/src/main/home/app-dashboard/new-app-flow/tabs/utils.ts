export const DATABASE_HEIGHT_ENABLED = 374;
export const DATABASE_HEIGHT_DISABLED = 119;
export const RESOURCE_HEIGHT_WITHOUT_AUTOSCALING = 446;
export const RESOURCE_HEIGHT_WITH_AUTOSCALING = 833;
export const MIB_TO_GIB = 1024;
export const MILI_TO_CORE = 1000;
interface InstanceDetails {
    vCPU: number;
    "RAM": number;
}

interface InstanceTypes {
    [key: string]: {
        [size: string]: InstanceDetails;
    };
}

export const AWS_INSTANCE_LIMITS: InstanceTypes = {
    "t3a": {
        "nano": { "vCPU": 2, "RAM": 0.5 },
        "micro": { "vCPU": 2, "RAM": 1 },
        "small": { "vCPU": 2, "RAM": 2 },
        "medium": { "vCPU": 2, "RAM": 4 },
        "large": { "vCPU": 2, "RAM": 8 },
        "xlarge": { "vCPU": 4, "RAM": 16 },
        "2xlarge": { "vCPU": 8, "RAM": 32 }
    },
    "t3": {
        "nano": { "vCPU": 2, "RAM": 0.5 },
        "micro": { "vCPU": 2, "RAM": 1 },
        "small": { "vCPU": 2, "RAM": 2 },
        "medium": { "vCPU": 2, "RAM": 4 },
        "large": { "vCPU": 2, "RAM": 8 },
        "xlarge": { "vCPU": 4, "RAM": 16 },
        "2xlarge": { "vCPU": 8, "RAM": 32 }
    },
    "t2": {
        "nano": { "vCPU": 1, "RAM": 0.5 },
        "micro": { "vCPU": 1, "RAM": 1 },
        "small": { "vCPU": 1, "RAM": 2 },
        "medium": { "vCPU": 2, "RAM": 4 },
        "large": { "vCPU": 2, "RAM": 8 },
        "xlarge": { "vCPU": 4, "RAM": 16 },
        "2xlarge": { "vCPU": 8, "RAM": 32 }
    },
    "c6i": {
        "large": { "vCPU": 2, "RAM": 4 },
        "xlarge": { "vCPU": 4, "RAM": 8 },
        "2xlarge": { "vCPU": 8, "RAM": 16 },
        "4xlarge": { "vCPU": 16, "RAM": 32 },
        "8xlarge": { "vCPU": 32, "RAM": 64 },
        "12xlarge": { "vCPU": 48, "RAM": 96 },
    },
    "g4dn": {
        "xlarge": { "vCPU": 4, "RAM": 16 },
        "2xlarge": { "vCPU": 8, "RAM": 32 },
        "4xlarge": { "vCPU": 16, "RAM": 64 },
        "8xlarge": { "vCPU": 32, "RAM": 128 },
    },
    "r6a": {
        "large": { "vCPU": 2, "RAM": 16 },
        "xlarge": { "vCPU": 4, "RAM": 32 },
        "2xlarge": { "vCPU": 8, "RAM": 64 },
        "4xlarge": { "vCPU": 16, "RAM": 128 },
        "8xlarge": { "vCPU": 32, "RAM": 256 },
    },
    "c5": {
        "large": { "vCPU": 2, "RAM": 4 },
        "xlarge": { "vCPU": 4, "RAM": 8 },
        "2xlarge": { "vCPU": 8, "RAM": 16 },
        "4xlarge": { "vCPU": 16, "RAM": 32 },
    },
    "m5": {
        "large": { "vCPU": 2, "RAM": 8 },
        "xlarge": { "vCPU": 4, "RAM": 16 },
        "2xlarge": { "vCPU": 8, "RAM": 32 },
        "4xlarge": { "vCPU": 16, "RAM": 64 },
    },
    "x2gd": {
        "medium": { "vCPU": 1, "RAM": 16 },
        "large": { "vCPU": 2, "RAM": 32 },
        "xlarge": { "vCPU": 4, "RAM": 64 },
    }
}


export const UPPER_BOUND_SMART = .5;
export const RESOURCE_ALLOCATION_RAM = 1;