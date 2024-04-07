/*
 * Copyright 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.springframework.cloud.dataflow.server.db.migration.mysql;

import org.springframework.cloud.dataflow.server.db.migration.AbstractBootVersionMigration;

/**
 * The boot_version field will provide for indicating the version of
 * Spring Boot used by the application and by implication the
 * schema version of task and batch tables.
 * @author Corneil du Plessis
 * @since 2.11
 */
public class V6__Boot3_Boot_Version extends AbstractBootVersionMigration {
}
