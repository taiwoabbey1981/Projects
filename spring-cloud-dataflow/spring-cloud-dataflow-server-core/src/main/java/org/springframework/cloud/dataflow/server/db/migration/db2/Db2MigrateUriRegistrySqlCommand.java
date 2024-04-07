/*
 * Copyright 2019 the original author or authors.
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

import java.util.List;

import org.postgresql.core.SqlCommand;

import org.springframework.cloud.dataflow.server.db.migration.AbstractMigrateUriRegistrySqlCommand;
import org.springframework.jdbc.core.JdbcTemplate;

/**
 * {@code db2} related {@link SqlCommand} for migrating data from
 * {@code URI_REGISTRY} into {@code app_registration}.
 *
 * @author Janne Valkealahti
 *
 */
public class Db2MigrateUriRegistrySqlCommand extends AbstractMigrateUriRegistrySqlCommand {

	@Override
	protected void updateAppRegistration(JdbcTemplate jdbcTemplate, List<AppRegistrationMigrationData> data) {
		for (AppRegistrationMigrationData d : data) {
			Long nextVal = jdbcTemplate.queryForObject("values next value for hibernate_sequence", Long.class);
			jdbcTemplate.update(
					"insert into app_registration (id, object_version, default_version, metadata_uri, name, type, uri, version) values (?,?,?,?,?,?,?,?)",
					nextVal, 0, d.isDefaultVersion() ? 1 : 0, d.getMetadataUri(), d.getName(), d.getType(), d.getUri(),
					d.getVersion());
		}
	}
}
