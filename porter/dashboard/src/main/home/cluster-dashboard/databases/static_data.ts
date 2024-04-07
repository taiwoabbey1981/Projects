export const POSTGRES_ENGINE_VERSIONS: {
  [key: string]: { label: string; value: string }[];
} = {
  postgres9: [
    { label: "v9.6.1", value: "9.6.1" },
    { label: "v9.6.2", value: "9.6.2" },
    { label: "v9.6.3", value: "9.6.3" },
    { label: "v9.6.4", value: "9.6.4" },
    { label: "v9.6.5", value: "9.6.5" },
    { label: "v9.6.6", value: "9.6.6" },
    { label: "v9.6.7", value: "9.6.7" },
    { label: "v9.6.8", value: "9.6.8" },
    { label: "v9.6.9", value: "9.6.9" },
    { label: "v9.6.10", value: "9.6.10" },
    { label: "v9.6.11", value: "9.6.11" },
    { label: "v9.6.12", value: "9.6.12" },
    { label: "v9.6.13", value: "9.6.13" },
    { label: "v9.6.14", value: "9.6.14" },
    { label: "v9.6.15", value: "9.6.15" },
    { label: "v9.6.16", value: "9.6.16" },
    { label: "v9.6.17", value: "9.6.17" },
    { label: "v9.6.18", value: "9.6.18" },
    { label: "v9.6.19", value: "9.6.19" },
    { label: "v9.6.20", value: "9.6.20" },
    { label: "v9.6.21", value: "9.6.21" },
    { label: "v9.6.22", value: "9.6.22" },
    { label: "v9.6.23", value: "9.6.23" },
  ],
  postgres10: [
    { label: "v10.1", value: "10.1" },
    { label: "v10.2", value: "10.2" },
    { label: "v10.3", value: "10.3" },
    { label: "v10.4", value: "10.4" },
    { label: "v10.5", value: "10.5" },
    { label: "v10.6", value: "10.6" },
    { label: "v10.7", value: "10.7" },
    { label: "v10.8", value: "10.8" },
    { label: "v10.9", value: "10.9" },
    { label: "v10.10", value: "10.10" },
    { label: "v10.11", value: "10.11" },
    { label: "v10.12", value: "10.12" },
    { label: "v10.13", value: "10.13" },
    { label: "v10.14", value: "10.14" },
    { label: "v10.15", value: "10.15" },
    { label: "v10.16", value: "10.16" },
    { label: "v10.17", value: "10.17" },
    { label: "v10.18", value: "10.18" },
  ],
  postgres11: [
    { label: "v11.1", value: "11.1" },
    { label: "v11.2", value: "11.2" },
    { label: "v11.3", value: "11.3" },
    { label: "v11.4", value: "11.4" },
    { label: "v11.5", value: "11.5" },
    { label: "v11.6", value: "11.6" },
    { label: "v11.7", value: "11.7" },
    { label: "v11.8", value: "11.8" },
    { label: "v11.9", value: "11.9" },
    { label: "v11.10", value: "11.10" },
    { label: "v11.11", value: "11.11" },
    { label: "v11.12", value: "11.12" },
    { label: "v11.13", value: "11.13" },
  ],
  postgres12: [
    { label: "v12.2", value: "12.2" },
    { label: "v12.3", value: "12.3" },
    { label: "v12.4", value: "12.4" },
    { label: "v12.5", value: "12.5" },
    { label: "v12.6", value: "12.6" },
    { label: "v12.7", value: "12.7" },
    { label: "v12.8", value: "12.8" },
  ],
  postgres13: [
    { label: "v13.1", value: "13.1" },
    { label: "v13.2", value: "13.2" },
    { label: "v13.3", value: "13.3" },
    { label: "v13.4", value: "13.4" },
  ],
};

export const POSTGRES_DB_FAMILIES = (() => {
  const dbFamilies = Object.keys(POSTGRES_ENGINE_VERSIONS);
  return dbFamilies.map((family) => {
    return {
      label: family,
      value: family,
    };
  });
})();

export const DEFAULT_DB_FAMILY =
  POSTGRES_DB_FAMILIES[POSTGRES_DB_FAMILIES.length - 1].value;

export const LAST_POSTGRES_ENGINE_VERSION =
  POSTGRES_ENGINE_VERSIONS.postgres13[
    POSTGRES_ENGINE_VERSIONS.postgres13.length - 1
  ].value;

export const FORM_DEFAULT_VALUES = {
  db_allocated_storage: "10",
  db_max_allocated_storage: "20",
  db_storage_encrypted: false,
};

export const DATABASE_INSTANCE_TYPES = [
  { value: "db.t2.medium", label: "db.t2.medium" },
  { value: "db.t2.xlarge", label: "db.t2.xlarge" },
  { value: "db.t2.2xlarge", label: "db.t2.2xlarge" },
  { value: "db.t3.medium", label: "db.t3.medium" },
  { value: "db.t3.xlarge", label: "db.t3.xlarge" },
  { value: "db.t3.2xlarge", label: "db.t3.2xlarge" },
];

export const DEFAULT_DATABASE_INSTANCE_TYPE = DATABASE_INSTANCE_TYPES[0].value;
