/*
 * Copyright 2018-2019 the original author or authors.
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
package org.springframework.cloud.dataflow.common.test.docker.compose.connection;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.not;
import static org.hamcrest.core.Is.is;
import static org.springframework.cloud.dataflow.common.test.docker.compose.configuration.DaemonHostIpResolver.LOCALHOST;
import static org.springframework.cloud.dataflow.common.test.docker.compose.configuration.DockerType.DAEMON;
import static org.springframework.cloud.dataflow.common.test.docker.compose.configuration.DockerType.REMOTE;
import static org.springframework.cloud.dataflow.common.test.docker.compose.configuration.EnvironmentVariables.DOCKER_CERT_PATH;
import static org.springframework.cloud.dataflow.common.test.docker.compose.configuration.EnvironmentVariables.DOCKER_HOST;
import static org.springframework.cloud.dataflow.common.test.docker.compose.configuration.EnvironmentVariables.DOCKER_TLS_VERIFY;
import static org.springframework.cloud.dataflow.common.test.docker.compose.matchers.DockerMachineEnvironmentMatcher.containsEnvironment;

import java.util.HashMap;
import java.util.Map;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.springframework.cloud.dataflow.common.test.docker.compose.connection.DockerMachine;
import org.springframework.cloud.dataflow.common.test.docker.compose.connection.DockerMachine.LocalBuilder;

public class LocalBuilderTests {

	@Rule
	public ExpectedException exception = ExpectedException.none();

	@Test
	public void override_previous_environment_when_additional_environment_set_twice_daemon() {
		Map<String, String> environment1 = new HashMap<>();
		environment1.put("ENV_1", "VAL_1");
		Map<String, String> environment2 = new HashMap<>();
		environment2.put("ENV_2", "VAL_2");
		DockerMachine localMachine = new LocalBuilder(DAEMON, new HashMap<>()).withEnvironment(environment1)
																		   .withEnvironment(environment2)
																		   .build();
		assertThat(localMachine, not(containsEnvironment(environment1)));
		assertThat(localMachine, containsEnvironment(environment2));
	}

	@Test
	public void be_union_of_additional_environment_and_individual_environment_when_both_set_daemon() {
		Map<String, String> environment = new HashMap<>();
		environment.put("ENV_1", "VAL_1");
		environment.put("ENV_2", "VAL_2");
		DockerMachine localMachine = new LocalBuilder(DAEMON, new HashMap<>()).withEnvironment(environment)
																		   .withAdditionalEnvironmentVariable("ENV_3", "VAL_3")
																		   .build();
		assertThat(localMachine, containsEnvironment(environment));
		Map<String, String> environment2 = new HashMap<>();
		environment2.put("ENV_3", "VAL_3");
		assertThat(localMachine, containsEnvironment(environment2));
	}

	@Test
	public void override_previous_environment_with_additional_environment_set_twice_remote() {
		Map<String, String> dockerVariables = new HashMap<>();
		dockerVariables.put(DOCKER_HOST, "tcp://192.168.99.100:2376");
		Map<String, String> environment1 = new HashMap<>();
		environment1.put("ENV_1", "VAL_1");
		Map<String, String> environment2 = new HashMap<>();
		environment2.put("ENV_2", "VAL_2");
		DockerMachine localMachine = new LocalBuilder(REMOTE, dockerVariables).withEnvironment(environment1)
																		   .withEnvironment(environment2)
																		   .build();
		assertThat(localMachine, not(containsEnvironment(environment1)));
		assertThat(localMachine, containsEnvironment(environment2));
	}

	@Test
	public void be_union_of_additional_environment_and_individual_environment_when_both_set_remote() {
		Map<String, String> dockerVariables = new HashMap<>();
		dockerVariables.put(DOCKER_HOST, "tcp://192.168.99.100:2376");
		Map<String, String> environment = new HashMap<>();
		environment.put("ENV_1", "VAL_1");
		environment.put("ENV_2", "VAL_2");
		DockerMachine localMachine = new LocalBuilder(REMOTE, dockerVariables).withEnvironment(environment)
																			  .withAdditionalEnvironmentVariable("ENV_3", "VAL_3")
																			  .build();
		assertThat(localMachine, containsEnvironment(environment));
		Map<String, String> environment2 = new HashMap<>();
		environment2.put("ENV_3", "VAL_3");
		assertThat(localMachine, containsEnvironment(environment2));
	}

	@Test
	public void get_variable_overriden_with_additional_environment() {
		Map<String, String> environment = new HashMap<>();
		environment.put("ENV_1", "VAL_1");
		environment.put("ENV_2", "VAL_2");
		DockerMachine localMachine = new LocalBuilder(DAEMON, new HashMap<>()).withEnvironment(environment)
																		   .withAdditionalEnvironmentVariable("ENV_2", "DIFFERENT_VALUE")
																		   .build();

		Map<String, String> expected = new HashMap<>();
		expected.put("ENV_1", "VAL_1");
		expected.put("ENV_2", "DIFFERENT_VALUE");
		assertThat(localMachine, not(containsEnvironment(environment)));
		assertThat(localMachine, containsEnvironment(expected));
	}

	@Test
	public void override_system_environment_with_additional_environment() {
		Map<String, String> systemEnv = new HashMap<>();
		systemEnv.put("ENV_1", "VAL_1");
		Map<String, String> overrideEnv = new HashMap<>();
		overrideEnv.put("ENV_1", "DIFFERENT_VALUE");
		DockerMachine localMachine = new LocalBuilder(DAEMON, systemEnv)
				.withEnvironment(overrideEnv)
				.build();

		assertThat(localMachine, not(containsEnvironment(systemEnv)));
		assertThat(localMachine, containsEnvironment(overrideEnv));
	}

	@Test
	public void have_invalid_variables_daemon() {
		Map<String, String> invalidDockerVariables = new HashMap<>();
		invalidDockerVariables.put(DOCKER_HOST, "tcp://192.168.99.100:2376");
		invalidDockerVariables.put(DOCKER_TLS_VERIFY, "1");
		invalidDockerVariables.put(DOCKER_CERT_PATH, "/path/to/certs");

		exception.expect(IllegalStateException.class);
		exception.expectMessage("These variables were set");
		exception.expectMessage(DOCKER_HOST);
		exception.expectMessage(DOCKER_CERT_PATH);
		exception.expectMessage(DOCKER_TLS_VERIFY);
		exception.expectMessage("They cannot be set when connecting to a local docker daemon");

		new LocalBuilder(DAEMON, invalidDockerVariables).build();
	}

	@Test
	public void have_invalid_additional_variables_daemon() {
		exception.expect(IllegalStateException.class);
		exception.expectMessage("The following variables");
		exception.expectMessage(DOCKER_HOST);
		exception.expectMessage("cannot exist in your additional environment variable block");

		new LocalBuilder(DAEMON, new HashMap<>()).withAdditionalEnvironmentVariable(DOCKER_HOST, "tcp://192.168.99.100:2376")
											  .build();
	}

	@Test
	public void have_invalid_additional_variables_remote() {
		Map<String, String> dockerVariables = new HashMap<>();
		dockerVariables.put(DOCKER_HOST, "tcp://192.168.99.100:2376");
		dockerVariables.put(DOCKER_TLS_VERIFY, "1");
		dockerVariables.put(DOCKER_CERT_PATH, "/path/to/certs");

		exception.expect(IllegalStateException.class);
		exception.expectMessage("The following variables");
		exception.expectMessage(DOCKER_HOST);
		exception.expectMessage("cannot exist in your additional environment variable block");

		new LocalBuilder(REMOTE, dockerVariables).withAdditionalEnvironmentVariable(DOCKER_HOST, "tcp://192.168.99.101:2376")
												 .build();
	}

	@Test
	public void return_localhost_as_ip_daemon() {
		DockerMachine localMachine = new LocalBuilder(DAEMON, new HashMap<>()).build();
		assertThat(localMachine.getIp(), is(LOCALHOST));
	}

	@Test
	public void return_docker_host_as_ip_remote() {
		Map<String, String> dockerVariables = new HashMap<>();
		dockerVariables.put(DOCKER_HOST, "tcp://192.168.99.100:2376");
		dockerVariables.put(DOCKER_TLS_VERIFY, "1");
		dockerVariables.put(DOCKER_CERT_PATH, "/path/to/certs");

		DockerMachine localMachine = new LocalBuilder(REMOTE, dockerVariables).build();
		assertThat(localMachine.getIp(), is("192.168.99.100"));
	}

	@Test
	public void have_missing_docker_host_remote() {
		exception.expect(IllegalStateException.class);
		exception.expectMessage("Missing required environment variables: ");
		exception.expectMessage(DOCKER_HOST);
		new LocalBuilder(REMOTE, new HashMap<>()).build();
	}

	@Test
	public void build_without_tls_remote() {
		Map<String, String> dockerVariables = new HashMap<>();
		dockerVariables.put(DOCKER_HOST, "tcp://192.168.99.100:2376");

		DockerMachine localMachine = new LocalBuilder(REMOTE, dockerVariables).build();
		assertThat(localMachine, containsEnvironment(dockerVariables));
	}

	@Test
	public void have_missing_cert_path_remote() {
		Map<String, String> dockerVariables = new HashMap<>();
		dockerVariables.put(DOCKER_HOST, "tcp://192.168.99.100:2376");
		dockerVariables.put(DOCKER_TLS_VERIFY, "1");

		exception.expect(IllegalStateException.class);
		exception.expectMessage("Missing required environment variables: ");
		exception.expectMessage(DOCKER_CERT_PATH);
		new LocalBuilder(REMOTE, dockerVariables).build();
	}

	@Test
	public void build_with_tls_remote() {
		Map<String, String> dockerVariables = new HashMap<>();
		dockerVariables.put(DOCKER_HOST, "tcp://192.168.99.100:2376");
		dockerVariables.put(DOCKER_TLS_VERIFY, "1");
		dockerVariables.put(DOCKER_CERT_PATH, "/path/to/certs");

		DockerMachine localMachine = new LocalBuilder(REMOTE, dockerVariables).build();
		assertThat(localMachine, containsEnvironment(dockerVariables));
	}
}
