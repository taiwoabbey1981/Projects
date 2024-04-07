/*
 * Copyright 2022-2022 the original author or authors.
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
package org.flywaydb.database.mysql.mariadb;

import org.flywaydb.core.api.configuration.Configuration;
import org.flywaydb.core.internal.jdbc.JdbcConnectionFactory;
import org.flywaydb.core.internal.jdbc.StatementInterceptor;
import org.flywaydb.database.mysql.MySQL57Database;
import org.flywaydb.database.mysql.mariadb.MariaDBDatabase;

public class MariaDB57Database extends MySQL57Database {

    public MariaDB57Database(Configuration configuration, JdbcConnectionFactory jdbcConnectionFactory, StatementInterceptor statementInterceptor) {
        super(configuration, jdbcConnectionFactory, statementInterceptor, new MariaDBDatabase(configuration, jdbcConnectionFactory, statementInterceptor));
    }
}
