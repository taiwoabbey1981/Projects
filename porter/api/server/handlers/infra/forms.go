package infra

const testForm = `name: Test
hasSource: false
includeHiddenFields: true
isClusterScoped: true
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: String to echo
    - type: string-input
      variable: echo
      settings:
        default: hello
`

const s3Form = `name: S3
hasSource: false
includeHiddenFields: true
isClusterScoped: true
tabs:
- name: main
  label: Main
  sections:
  - name: heading
    contents:
    - type: heading
      label: S3 Settings
  - name: bucket_name
    contents:
    - type: string-input
      label: Bucket Name
      required: true
      placeholder: "s3-bucket-name"
      variable: bucket_name
`

const rdsForm = `name: RDS
hasSource: false
includeHiddenFields: true
isClusterScoped: true
tabs:
- name: main
  label: Main
  sections:
  - name: heading
    contents:
    - type: heading
      label: Database Settings
  - name: user
    contents:
    - type: string-input
      label: Database Master User
      required: true
      placeholder: "admin"
      variable: db_user
  - name: password
    contents:
    - type: string-input
      required: true
      label: Database Master Password
      variable: db_passwd
  - name: name
    contents:
    - type: string-input
      label: Database Name
      required: true
      placeholder: "rds-staging"
      variable: db_name
  - name: machine-type
    contents:
    - type: select
      label: ⚙️ Database Machine Type
      variable: machine_type
      settings:
        default: db.t3.medium
        options:
        - label: db.t2.medium
          value: db.t2.medium
        - label: db.t2.xlarge
          value: db.t2.xlarge
        - label: db.t2.2xlarge
          value: db.t2.2xlarge
        - label: db.t3.medium
          value: db.t3.medium
        - label: db.t3.xlarge
          value: db.t3.xlarge
        - label: db.t3.2xlarge
          value: db.t3.2xlarge
        - label: db.r5.large
          value: db.r5.large
        - label: db.r5.xlarge
          value: db.r5.xlarge
        - label: db.r5.2xlarge
          value: db.r5.2xlarge
  - name: family-versions
    contents:
    - type: select
      label:  Database Family Version
      variable: db_family
      settings:
        default: postgres13
        options:
        - label: "Postgres 9"
          value: postgres9
        - label: "Postgres 10"
          value: postgres10
        - label: "Postgres 11"
          value: postgres11
        - label: "Postgres 12"
          value: postgres12
        - label: "Postgres 13"
          value: postgres13
        - label: "Postgres 14"
          value: postgres14
  - name: pg-9-versions
    show_if:
      is: "postgres9"
      variable: db_family
    contents:
    - type: select
      label:  Database Version
      variable: db_engine_version
      settings:
        default: "9.6.23"
        options:
        - label: "v9.6.1"
          value: "9.6.1"
        - label: "v9.6.2"
          value: "9.6.2"
        - label: "v9.6.3"
          value: "9.6.3"
        - label: "v9.6.4"
          value: "9.6.4"
        - label: "v9.6.5"
          value: "9.6.5"
        - label: "v9.6.6"
          value: "9.6.6"
        - label: "v9.6.7"
          value: "9.6.7"
        - label: "v9.6.8"
          value: "9.6.8"
        - label: "v9.6.10"
          value: "9.6.10"
        - label: "v9.6.11"
          value: "9.6.11"
        - label: "v9.6.12"
          value: "9.6.12"
        - label: "v9.6.13"
          value: "9.6.13"
        - label: "v9.6.14"
          value: "9.6.14"
        - label: "v9.6.15"
          value: "9.6.15"
        - label: "v9.6.16"
          value: "9.6.16"
        - label: "v9.6.17"
          value: "9.6.17"
        - label: "v9.6.18"
          value: "9.6.18"
        - label: "v9.6.19"
          value: "9.6.19"
        - label: "v9.6.20"
          value: "9.6.20"
        - label: "v9.6.21"
          value: "9.6.21"
        - label: "v9.6.22"
          value: "9.6.22"
        - label: "v9.6.23"
          value: "9.6.23"
  - name: pg-10-versions
    show_if:
      is: "postgres10"
      variable: db_family
    contents:
    - type: select
      label:  Database Version
      variable: db_engine_version
      settings:
        default: "10.22"
        options:
        - label: "v10.1"
          value: "10.1"
        - label: "v10.2"
          value: "10.2"
        - label: "v10.3"
          value: "10.3"
        - label: "v10.4"
          value: "10.4"
        - label: "v10.5"
          value: "10.5"
        - label: "v10.6"
          value: "10.6"
        - label: "v10.7"
          value: "10.7"
        - label: "v10.8"
          value: "10.8"
        - label: "v10.9"
          value: "10.9"
        - label: "v10.10"
          value: "10.10"
        - label: "v10.11"
          value: "10.11"
        - label: "v10.12"
          value: "10.12"
        - label: "v10.13"
          value: "10.13"
        - label: "v10.14"
          value: "10.14"
        - label: "v10.15"
          value: "10.15"
        - label: "v10.16"
          value: "10.16"
        - label: "v10.17"
          value: "10.17"
        - label: "v10.18"
          value: "10.18"
        - label: "v10.19"
          value: "10.19"
        - label: "v10.20"
          value: "10.20"
        - label: "v10.21"
          value: "10.21"
        - label: "v10.22"
          value: "10.22"
  - name: pg-11-versions
    show_if:
      is: "postgres11"
      variable: db_family
    contents:
    - type: select
      label:  Database Version
      variable: db_engine_version
      settings:
        default: "11.17"
        options:
        - label: "v11.1"
          value: "11.1"
        - label: "v11.2"
          value: "11.2"
        - label: "v11.3"
          value: "11.3"
        - label: "v11.4"
          value: "11.4"
        - label: "v11.5"
          value: "11.5"
        - label: "v11.6"
          value: "11.6"
        - label: "v11.7"
          value: "11.7"
        - label: "v11.8"
          value: "11.8"
        - label: "v11.9"
          value: "11.9"
        - label: "v11.10"
          value: "11.10"
        - label: "v11.11"
          value: "11.11"
        - label: "v11.12"
          value: "11.12"
        - label: "v11.13"
          value: "11.13"
        - label: "v11.14"
          value: "11.14"
        - label: "v11.15"
          value: "11.15"
        - label: "v11.16"
          value: "11.16"
        - label: "v11.17"
          value: "11.17"
  - name: pg-12-versions
    show_if:
      is: "postgres12"
      variable: db_family
    contents:
    - type: select
      label:  Database Version
      variable: db_engine_version
      settings:
        default: "12.12"
        options:
        - label: "v12.2"
          value: "12.2"
        - label: "v12.3"
          value: "12.3"
        - label: "v12.4"
          value: "12.4"
        - label: "v12.5"
          value: "12.5"
        - label: "v12.6"
          value: "12.6"
        - label: "v12.7"
          value: "12.7"
        - label: "v12.8"
          value: "12.8"
        - label: "v12.9"
          value: "12.9"
        - label: "v12.10"
          value: "12.10"
        - label: "v12.11"
          value: "12.11"
        - label: "v12.12"
          value: "12.12"
  - name: pg-13-versions
    show_if:
      is: "postgres13"
      variable: db_family
    contents:
    - type: select
      label:  Database Version
      variable: db_engine_version
      settings:
        default: "13.8"
        options:
        - label: "v13.1"
          value: "13.1"
        - label: "v13.2"
          value: "13.2"
        - label: "v13.3"
          value: "13.3"
        - label: "v13.4"
          value: "13.4"
        - label: "v13.5"
          value: "13.5"
        - label: "v13.6"
          value: "13.6"
        - label: "v13.7"
          value: "13.7"
        - label: "v13.8"
          value: "13.8"
  - name: pg-14-versions
    show_if:
      is: "postgres14"
      variable: db_family
    contents:
    - type: select
      label:  Database Version
      variable: db_engine_version
      settings:
        default: "14.5"
        options:
        - label: "v14.1"
          value: "14.1"
        - label: "v14.2"
          value: "14.2"
        - label: "v14.3"
          value: "14.3"
        - label: "v14.4"
          value: "14.4"
        - label: "v14.5"
          value: "14.5"
  - name: additional-settings
    contents:
    - type: heading
      label: Additional Settings
    - type: checkbox
      variable: db_deletion_protection
      label: Enable deletion protection for the database.
      settings:
        default: false
- name: storage
  label: Storage
  sections:
  - name: storage
    contents:
    - type: heading
      label: Storage Settings
    - type: number-input
      label: Specify the amount of storage to allocate to this instance in gigabytes.
      variable: db_allocated_storage
      placeholder: "ex: 10"
      settings:
        default: 10
    - type: number-input
      label: Specify the maximum storage that this instance can scale to in gigabytes.
      variable: db_max_allocated_storage
      placeholder: "ex: 20"
      settings:
        default: 20
    - type: checkbox
      variable: db_storage_encrypted
      label: Enable storage encryption for the database.
      settings:
        default: false
- name: advanced
  label: Advanced
  sections:
  - name: replicas
    contents:
    - type: heading
      label: Read Replicas
    - type: subtitle
      label: Specify the number of read replicas to run alongside your RDS instance.
    - type: number-input
      label: Replicas
      variable: db_replicas
      placeholder: "ex: 1"
      settings:
        default: 0`

const ecrForm = `name: ECR
hasSource: false
includeHiddenFields: true
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: ECR Configuration
    - type: string-input
      label: ECR Name
      required: true
      placeholder: my-awesome-registry
      variable: ecr_name
`

const eksForm = `name: EKS
hasSource: false
includeHiddenFields: true
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: EKS Configuration
    - type: select
      label: ⚙️ AWS Machine Type
      variable: machine_type
      settings:
        default: t2.medium
        options:
        - label: t2.medium
          value: t2.medium
        - label: t2.large
          value: t2.large
        - label: t2.xlarge
          value: t2.xlarge
        - label: t2.2xlarge
          value: t2.2xlarge
        - label: t3.medium
          value: t3.medium
        - label: t3.large
          value: t3.large
        - label: t3.xlarge
          value: t3.xlarge
        - label: t3.2xlarge
          value: t3.2xlarge
        - label: c5.large
          value: c5.large
        - label: c5.xlarge
          value: c5.xlarge
        - label: c5.2xlarge
          value: c5.2xlarge
        - label: c6a.large
          value: c6a.large
        - label: c6a.xlarge
          value: c6a.xlarge
        - label: c6a.2xlarge
          value: c6a.2xlarge
        - label: c6a.4xlarge
          value: c6a.4xlarge
        - label: c6i.large
          value: c6i.large
        - label: c6i.xlarge
          value: c6i.xlarge
        - label: c6i.2xlarge
          value: c6i.2xlarge
        - label: c6i.4xlarge
          value: c6i.4xlarge
        - label: g4dn.xlarge
          value: g4dn.xlarge
        - label: g4dn.2xlarge
          value: g4dn.2xlarge
        - label: g4dn.4xlarge
          value: g4dn.4xlarge
        - label: g5.xlarge
          value: g5.xlarge
        - label: g5.2xlarge
          value: g5.2xlarge
        - label: g5.4xlarge
          value: g5.4xlarge
        - label: m6a.large
          value: m6a.large
        - label: m6a.xlarge
          value: m6a.xlarge
        - label: m6a.2xlarge
          value: m6a.2xlarge
        - label: m6a.4xlarge
          value: m6a.4xlarge
        - label: m6i.large
          value: m6i.large
        - label: m6i.xlarge
          value: m6i.xlarge
        - label: m6i.2xlarge
          value: m6i.2xlarge
        - label: m6i.4xlarge
          value: m6i.4xlarge
        - label: r5.large
          value: r5.large
        - label: r5.xlarge
          value: r5.xlarge
    - type: string-input
      label: 👤 Issuer Email
      required: true
      placeholder: example@example.com
      variable: issuer_email
    - type: string-input
      label: EKS Cluster Name
      required: true
      placeholder: my-cluster
      variable: cluster_name
    - type: select
      label: EKS control plane version
      variable: cluster_version
      settings:
        default: "1.24"
        options:
        - label: "1.20"
          value: "1.20"
        - label: "1.21"
          value: "1.21"
        - label: "1.22"
          value: "1.22"
        - label: "1.23"
          value: "1.23"
        - label: "1.24"
          value: "1.24"
    - type: number-input
      label: Minimum number of EC2 instances to create in the application autoscaling group.
      variable: min_instances
      placeholder: "ex: 1"
      settings:
        default: 1
    - type: number-input
      label: Maximum number of EC2 instances to create in the application autoscaling group.
      variable: max_instances
      placeholder: "ex: 10"
      settings:
        default: 10
- name: additional_nodegroup
  label: Additional Node Groups
  sections:
  - name: is_additional_enabled
    contents:
    - type: heading
      label: Additional Node Groups
    - type: checkbox
      variable: additional_nodegroup_enabled
      label: Enable an additional node group for this cluster.
      settings:
        default: false
  - name: additional_settings
    show_if: additional_nodegroup_enabled
    contents:
    - type: string-input
      label: Label for this node group. Multiple labels should be comma separated.
      variable: additional_nodegroup_label
      placeholder: "ex: porter.run/workload-kind=job"
      settings:
        default: porter.run/workload-kind=database
    - type: string-input
      label: Taint for this node group.
      variable: additional_nodegroup_taint
      placeholder: "ex: porter.run/workload-kind=job:NoSchedule"
      settings:
        default: porter.run/workload-kind=database:NoSchedule
    - type: checkbox
      variable: additional_stateful_nodegroup_enabled
      label: Stateful Workload
      settings:
        default: false
    - type: select
      label: ⚙️ AWS System Machine Type
      variable: additional_nodegroup_machine_type
      settings:
        default: t2.medium
        options:
        - label: t2.medium
          value: t2.medium
        - label: t2.large
          value: t2.large
        - label: t2.xlarge
          value: t2.xlarge
        - label: t2.2xlarge
          value: t2.2xlarge
        - label: t3.medium
          value: t3.medium
        - label: t3.large
          value: t3.large
        - label: t3.xlarge
          value: t3.xlarge
        - label: t3.2xlarge
          value: t3.2xlarge
        - label: c6i.large
          value: c6i.large
        - label: c6i.xlarge
          value: c6i.xlarge
        - label: c6i.2xlarge
          value: c6i.2xlarge
        - label: c6i.4xlarge
          value: c6i.4xlarge
        - label: m6i.large
          value: m6i.large
        - label: m6i.xlarge
          value: m6i.xlarge
        - label: m6i.2xlarge
          value: m6i.2xlarge
        - label: m6i.4xlarge
          value: m6i.4xlarge
    - type: number-input
      label: Minimum number of EC2 instances to create in the application autoscaling group.
      variable: additional_nodegroup_min_instances
      placeholder: "ex: 1"
      settings:
        default: 1
    - type: number-input
      label: Maximum number of EC2 instances to create in the application autoscaling group.
      variable: additional_nodegroup_max_instances
      placeholder: "ex: 10"
      settings:
        default: 10
- name: iam
  label: IAM
  sections:
  - name: toggle_aws_auth
    contents:
    - type: heading
      label: Configure IAM Access
    - type: checkbox
      variable: manage_aws_auth_configmap
      label: Allow Porter to manage AWS authentication for the cluster.
      settings:
        default: true
  - name: aws_auth_warning
    show_if:
      not: manage_aws_auth_configmap
    contents:
    - type: subtitle
      label: "WARNING - turning this value off will result in the aws-auth configmap getting removed from the cluster, and will take existing AWS nodes offline until the configmap is re-added with the node's IAM role ARN. Make sure you know what you are doing."
  - name: arns
    show_if: manage_aws_auth_configmap
    contents:
    - type: heading
      label: Users
    - type: subtitle
      label: "Add AWS users to the cluster. The left input should be a valid AWS user ARN, and the right side should be a group on the cluster. For example, arn:aws:iam::66666666666:user/user1: system:masters."
    - type: key-value-array
      variable: aws_auth_users
      settings:
        default: {}
    - type: heading
      label: Roles
    - type: subtitle
      label: "Add AWS roles to the cluster. The left input should be a valid AWS role ARN, and the right side should be a group on the cluster. For example, arn:aws:iam::66666666666:role/role1: system:masters."
    - type: key-value-array
      variable: aws_auth_roles
      settings:
        default: {}
- name: advanced
  label: Advanced
  sections:
  - name: workload_machine_label
    contents:
    - type: heading
      label: Application Node Group Settings
    - type: string-input
      label: Add custom node labels to the application node group. If you are adding multiple labels, they should be comma separated.
      variable: custom_node_labels_application
      placeholder: "ex: mylabel=custom-label,mylabel2=another-one"
      settings:
        default: ""
  - name: system_machine_type
    contents:
    - type: heading
      label: System Node Group Settings
    - type: select
      label: ⚙️ AWS System Machine Type
      variable: system_machine_type
      settings:
        default: t2.medium
        options:
        - label: t2.medium
          value: t2.medium
        - label: t2.large
          value: t2.large
        - label: t2.xlarge
          value: t2.xlarge
        - label: t2.2xlarge
          value: t2.2xlarge
        - label: t3.medium
          value: t3.medium
        - label: t3.large
          value: t3.large
        - label: t3.xlarge
          value: t3.xlarge
        - label: t3.2xlarge
          value: t3.2xlarge
        - label: c6i.large
          value: c6i.large
        - label: c6i.xlarge
          value: c6i.xlarge
        - label: c6i.2xlarge
          value: c6i.2xlarge
        - label: c6i.4xlarge
          value: c6i.4xlarge
    - type: string-input
      label: Add custom node labels to the system node group. If you are adding multiple labels, they should be comma separated.
      variable: custom_node_labels_system
      placeholder: "ex: mylabel=custom-label,mylabel2=another-one"
      settings:
        default: ""
  - name: spot_instance_should_enable
    contents:
    - type: heading
      label: Spot Instance Settings
    - type: checkbox
      variable: spot_instances_enabled
      label: Enable spot instances for this cluster.
      settings:
        default: false
  - name: spot_instance_price
    show_if: spot_instances_enabled
    contents:
    - type: string-input
      label: Assign a bid price for the spot instance (optional).
      variable: spot_price
      placeholder: "ex: 0.05"
  - name: net_settings
    contents:
    - type: heading
      label: Networking Settings
    - type: string-input
      label: "Add a different CIDR range prefix (first two octets: for example 10.99 will create a VPC with CIDR range 10.99.0.0/16)."
      variable: cluster_vpc_cidr_octets
      placeholder: "ex: 10.99"
      settings:
        default: "10.99"
    - type: string-input
      label: "Add a different CIDR range prefix for cluster services (first two octets: for example 172.20 will configure EKS with CIDR range 172.20.0.0/16)."
      variable: cluster_service_cidr_octets
      placeholder: "ex: 172.20"
      settings:
        default: "172.20"
    - type: checkbox
      label: "Add additional private subnets to the cluster in each AZ."
      variable: additional_private_subnets
      settings:
        default: false
  - name: subnet_multiplicity
    show_if: additional_private_subnets
    contents:
    - type: number-input
      label: "Multiplicity of the subnet within each AZ."
      variable: additional_private_subnets_multiplicity
      settings:
        default: 3
  - name: net_settings_azs_toggle
    contents:
    - type: checkbox
      label: "Specify the AZs to provision this cluster in."
      variable: specify_azs
      settings:
        default: false
  - name: net_settings_azs
    show_if: specify_azs
    contents:
    - type: array-input
      variable: azs
      label: Availability Zones
  - name: net_settings_single_az_nat_gateway
    contents:
    - type: checkbox
      variable: single_az_nat_gateway
      label: "Place a NAT gateway inside a single AZ. Disabling this will place a NAT gateway in each AZ, for each subnet in your cluster's VPC."
      settings:
        default: true
  - name: nginx_settings
    contents:
    - type: heading
      label: NGINX Settings
    - type: checkbox
      variable: disable_nginx_load_balancer
      label: Disable NGINX load balancer and expose NGINX only on a cluster IP address.
      settings:
        default: false
  - name: prometheus_settings
    contents:
    - type: heading
      label: Prometheus Settings
    - type: checkbox
      variable: additional_prometheus_node_group
      label: Add an additional prometheus node group to ensure monitoring stability.
      settings:
        default: true
  - name: prometheus_machine_settings
    show_if: additional_prometheus_node_group
    contents:
    - type: select
      label: ⚙️ AWS Prometheus Machine Type
      variable: additional_prometheus_machine_type
      settings:
        default: t2.medium
        options:
        - label: t2.medium
          value: t2.medium
        - label: t2.large
          value: t2.large
        - label: t2.xlarge
          value: t2.xlarge
        - label: t3.medium
          value: t3.medium
        - label: t3.large
          value: t3.large
        - label: t3.xlarge
          value: t3.xlarge
    - type: string-input
      label: Add custom node labels to the monitoring node group. If you are adding multiple labels, they should be comma separated.
      variable: custom_node_labels_monitoring
      placeholder: "ex: mylabel=custom-label,mylabel2=another-one"
      settings:
        default: ""
  - name: kms_secret_encryption
    contents:
    - type: heading
      label: KMS Encryption
    - type: checkbox
      variable: is_kms_enabled
      label: Encrypt all Kubernetes secrets with AWS Key Management Service (KMS)
      settings:
        default: false
  - name: enable_logging
    contents:
    - type: heading
      label: Enable AWS Logging 
    - type: checkbox
      variable: vpc_flow_logs_enabled
      label: Enable VPC Flow Logs
      settings:
        default: false
    - type: checkbox
      variable: eks_control_plane_logs_enabled
      label: Enable EKS Control Plane Logs
      settings:
        default: false
  - name: add_custom_tags
    contents:
    - type: heading
      label: Add Custom Tags on AWS Resources Provisioned by Porter
    - type: key-value-array
      variable: custom_tags
      settings:
        default: {} 
`

const gcrForm = `name: GCR
hasSource: false
includeHiddenFields: true
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: GCR Configuration
    - type: select
      label: 📍 GCP Region
      variable: gcp_region
      settings:
        default: us-central1
        options:
        - label: asia-east1
          value: asia-east1
        - label: asia-east2
          value: asia-east2
        - label: asia-northeast1
          value: asia-northeast1
        - label: asia-northeast2
          value: asia-northeast2
        - label: asia-northeast3
          value: asia-northeast3
        - label: asia-south1
          value: asia-south1
        - label: asia-south2
          value: asia-south2
        - label: asia-southeast1
          value: asia-southeast1
        - label: asia-southeast2
          value: asia-southeast2
        - label: australia-southeast1
          value: australia-southeast1
        - label: australia-southeast2
          value: australia-southeast2
        - label: europe-central2
          value: europe-central2
        - label: europe-north1
          value: europe-north1
        - label: europe-west1
          value: europe-west1
        - label: europe-west2
          value: europe-west2
        - label: europe-west3
          value: europe-west3
        - label: europe-west4
          value: europe-west4
        - label: europe-west6
          value: europe-west6
        - label: europe-west8
          value: europe-west8
        - label: europe-west9
          value: europe-west9
        - label: europe-southwest1
          value: europe-southwest1
        - label: northamerica-northeast1
          value: northamerica-northeast1
        - label: northamerica-northeast2
          value: northamerica-northeast2
        - label: southamerica-east1
          value: southamerica-east1
        - label: southamerica-west1
          value: southamerica-west1
        - label: us-central1
          value: us-central1
        - label: us-east1
          value: us-east1
        - label: us-east4
          value: us-east4
        - label: us-east5
          value: us-east5
        - label: us-south1
          value: us-south1
        - label: us-west1
          value: us-west1
        - label: us-west2
          value: us-west2
        - label: us-west3
          value: us-west3
        - label: us-west4
          value: us-west4
`

const garForm = `name: GAR
hasSource: false
includeHiddenFields: true
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: GAR Configuration
    - type: select
      label: 📍 GCP Region
      variable: gcp_region
      settings:
        default: us-central1
        options:
        - label: asia-east1
          value: asia-east1
        - label: asia-east2
          value: asia-east2
        - label: asia-northeast1
          value: asia-northeast1
        - label: asia-northeast2
          value: asia-northeast2
        - label: asia-northeast3
          value: asia-northeast3
        - label: asia-south1
          value: asia-south1
        - label: asia-south2
          value: asia-south2
        - label: asia-southeast1
          value: asia-southeast1
        - label: asia-southeast2
          value: asia-southeast2
        - label: australia-southeast1
          value: australia-southeast1
        - label: australia-southeast2
          value: australia-southeast2
        - label: europe-central2
          value: europe-central2
        - label: europe-north1
          value: europe-north1
        - label: europe-west1
          value: europe-west1
        - label: europe-west2
          value: europe-west2
        - label: europe-west3
          value: europe-west3
        - label: europe-west4
          value: europe-west4
        - label: europe-west6
          value: europe-west6
        - label: europe-west8
          value: europe-west8
        - label: europe-west9
          value: europe-west9
        - label: europe-southwest1
          value: europe-southwest1
        - label: northamerica-northeast1
          value: northamerica-northeast1
        - label: northamerica-northeast2
          value: northamerica-northeast2
        - label: southamerica-east1
          value: southamerica-east1
        - label: southamerica-west1
          value: southamerica-west1
        - label: us-central1
          value: us-central1
        - label: us-east1
          value: us-east1
        - label: us-east4
          value: us-east4
        - label: us-east5
          value: us-east5
        - label: us-south1
          value: us-south1
        - label: us-west1
          value: us-west1
        - label: us-west2
          value: us-west2
        - label: us-west3
          value: us-west3
        - label: us-west4
          value: us-west4
        - label: us (multi-region)
          value: us
        - label: europe (multi-region)
          value: europe
        - label: asia (multi-region)
          value: asia
`

const gkeForm = `name: GKE
hasSource: false
includeHiddenFields: true
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: GKE Configuration
    - type: select
      label: 📍 GCP Region
      variable: gcp_region
      settings:
        default: us-central1
        options:
        - label: asia-east1
          value: asia-east1
        - label: asia-east2
          value: asia-east2
        - label: asia-northeast1
          value: asia-northeast1
        - label: asia-northeast2
          value: asia-northeast2
        - label: asia-northeast3
          value: asia-northeast3
        - label: asia-south1
          value: asia-south1
        - label: asia-south2
          value: asia-south2
        - label: asia-southeast1
          value: asia-southeast1
        - label: asia-southeast2
          value: asia-southeast2
        - label: australia-southeast1
          value: australia-southeast1
        - label: australia-southeast2
          value: australia-southeast2
        - label: europe-central2
          value: europe-central2
        - label: europe-north1
          value: europe-north1
        - label: europe-west1
          value: europe-west1
        - label: europe-west2
          value: europe-west2
        - label: europe-west3
          value: europe-west3
        - label: europe-west4
          value: europe-west4
        - label: europe-west6
          value: europe-west6
        - label: europe-west8
          value: europe-west8
        - label: europe-west9
          value: europe-west9
        - label: europe-southwest1
          value: europe-southwest1
        - label: northamerica-northeast1
          value: northamerica-northeast1
        - label: northamerica-northeast2
          value: northamerica-northeast2
        - label: southamerica-east1
          value: southamerica-east1
        - label: southamerica-west1
          value: southamerica-west1
        - label: us-central1
          value: us-central1
        - label: us-east1
          value: us-east1
        - label: us-east4
          value: us-east4
        - label: us-east5
          value: us-east5
        - label: us-south1
          value: us-south1
        - label: us-west1
          value: us-west1
        - label: us-west2
          value: us-west2
        - label: us-west3
          value: us-west3
        - label: us-west4
          value: us-west4
    - type: string-input
      label: 👤 Issuer Email
      required: true
      placeholder: example@example.com
      variable: issuer_email
    - type: string-input
      label: GKE Cluster Name
      required: true
      placeholder: my-cluster
      variable: cluster_name
`

const docrForm = `name: DOCR
hasSource: false
includeHiddenFields: true
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: DOCR Configuration
    - type: select
      label: DO Subscription Tier
      variable: docr_subscription_tier
      settings:
        default: basic
        options:
        - label: Basic
          value: basic
        - label: Professional
          value: professional
    - type: string-input
      label: DOCR Name
      required: true
      placeholder: my-awesome-registry
      variable: docr_name
`

const doksForm = `name: DOKS
hasSource: false
includeHiddenFields: true
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: DOKS Configuration
    - type: select
      label: 📍 DO Region
      variable: do_region
      settings:
        default: nyc1
        options:
        - label: Amsterdam 3
          value: ams3
        - label: Bangalore 1
          value: blr1
        - label: Frankfurt 1
          value: fra1
        - label: London 1
          value: lon1
        - label: New York 1
          value: nyc1
        - label: New York 3
          value: nyc3
        - label: San Francisco 2
          value: sfo2
        - label: San Francisco 3
          value: sfo3
        - label: Singapore 1
          value: sgp1
        - label: Toronto 1
          value: tor1
    - type: string-input
      label: 👤 Issuer Email
      required: true
      placeholder: example@example.com
      variable: issuer_email
    - type: string-input
      label: DOKS Cluster Name
      required: true
      placeholder: my-cluster
      variable: cluster_name
`

const acrForm = `name: ACR
hasSource: false
includeHiddenFields: true
isClusterScoped: false
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: ACR Configuration
    - type: select
      label: 📍 Azure Region
      variable: aks_region
      settings:
        default: East US
        options:
        - label: East US
          value: East US
        - label: East US 2
          value: East US 2
        - label: West US 2
          value: West US 2
        - label: West US 3
          value: West US 3
        - label: Norway East
          value: Norway East
    - type: string-input
      label: ACR Name
      required: true
      placeholder: my-registry
      variable: acr_name
`

const aksForm = `name: AKS
hasSource: false
includeHiddenFields: true
isClusterScoped: false
tabs:
- name: main
  label: Configuration
  sections:
  - name: section_one
    contents:
    - type: heading
      label: AKS Configuration
    - type: select
      label: 📍 Azure Region
      variable: aks_region
      settings:
        default: East US
        options:
        - label: East US
          value: East US
        - label: East US 2
          value: East US 2
        - label: West US 2
          value: West US 2
        - label: West US 3
          value: West US 3
        - label: Norway East
          value: Norway East
    - type: select
      label: ⚙️ Application Machine Type
      variable: app_machine_type
      settings:
        default: Standard_A2_v2
        options:
        - label: Standard A2
          value: Standard_A2_v2
        - label: Standard A4
          value: Standard_A4_v2
        - label: Standard D2
          value: Standard_D2_v3
    - type: string-input
      label: 👤 Issuer Email
      required: true
      placeholder: example@example.com
      variable: issuer_email
    - type: string-input
      label: AKS Cluster Name
      required: true
      placeholder: my-cluster
      variable: cluster_name
`
