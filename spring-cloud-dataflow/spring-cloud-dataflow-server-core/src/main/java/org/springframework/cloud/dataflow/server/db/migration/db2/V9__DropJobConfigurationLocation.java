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
package org.springframework.cloud.dataflow.server.db.migration.db2;

import org.springframework.cloud.dataflow.common.flyway.AbstractMigration;
import org.springframework.cloud.dataflow.server.db.migration.DropColumnSqlCommands;

import java.util.Collections;

/**
 * Removes extra JOB_CONFIGURATION_LOCATION columns.
 * @author Corneil du Plessis
 */
public class V9__DropJobConfigurationLocation extends AbstractMigration {
	public V9__DropJobConfigurationLocation() {
		super(Collections.singletonList(new DropColumnSqlCommands("BOOT3_BATCH_JOB_EXECUTION.JOB_CONFIGURATION_LOCATION")));
	}
}
