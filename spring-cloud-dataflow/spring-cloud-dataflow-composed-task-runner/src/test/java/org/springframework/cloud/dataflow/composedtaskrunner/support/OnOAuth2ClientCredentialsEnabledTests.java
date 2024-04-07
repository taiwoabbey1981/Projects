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
package org.springframework.cloud.dataflow.composedtaskrunner.support;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

import org.springframework.boot.test.util.TestPropertyValues;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Conditional;
import org.springframework.context.annotation.Configuration;

import static org.hamcrest.Matchers.equalTo;

/**
 * @author Gunnar Hillert
 */
public class OnOAuth2ClientCredentialsEnabledTests {

	private AnnotationConfigApplicationContext context;

	@AfterEach
	public void teardown() {
		if (this.context != null) {
			this.context.close();
		}
	}

	@Test
	public void noPropertySet() throws Exception {
		this.context = load(Config.class);
		org.hamcrest.MatcherAssert.assertThat(context.containsBean("myBean"), equalTo(false));
	}

	@Test
	public void propertyClientId() throws Exception {
		this.context = load(Config.class, "oauth2-client-credentials-client-id:12345");
		org.hamcrest.MatcherAssert.assertThat(context.containsBean("myBean"), equalTo(true));
	}

	@Test
	public void clientIdOnlyWithNoValue() throws Exception {
		this.context = load(Config.class, "oauth2-client-credentials-client-id:");
		org.hamcrest.MatcherAssert.assertThat(context.containsBean("myBean"), equalTo(false));
	}

	private AnnotationConfigApplicationContext load(Class<?> config, String... env) {
		AnnotationConfigApplicationContext context = new AnnotationConfigApplicationContext();
		TestPropertyValues.of(env).applyTo(context);
		context.register(config);
		context.refresh();
		return context;
	}

	@Configuration
	@Conditional(OnOAuth2ClientCredentialsEnabled.class)
	public static class Config {
		@Bean
		public String myBean() {
			return "myBean";
		}
	}
}
