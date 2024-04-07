/*
 * Copyright 2020-2020 the original author or authors.
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
package org.springframework.cloud.dataflow.container.registry;

import org.springframework.http.HttpHeaders;
import org.springframework.web.client.RestTemplate;

/**
 * @author Christian Tzolov
 */
public class ContainerRegistryRequest {

	private ContainerImage containerImage;
	private final ContainerRegistryConfiguration registryConf;
	private final HttpHeaders authHttpHeaders;
	private final RestTemplate restTemplate;

	public ContainerRegistryRequest(ContainerImage containerImage,
			ContainerRegistryConfiguration registryConf, HttpHeaders authHttpHeaders, RestTemplate requestRestTemplate) {
		this.containerImage = containerImage;
		this.registryConf = registryConf;
		this.authHttpHeaders = authHttpHeaders;
		this.restTemplate = requestRestTemplate;
	}

	public ContainerRegistryRequest(
			ContainerRegistryConfiguration registryConf, HttpHeaders authHttpHeaders, RestTemplate requestRestTemplate) {
		this.registryConf = registryConf;
		this.authHttpHeaders = authHttpHeaders;
		this.restTemplate = requestRestTemplate;
	}

	public ContainerImage getContainerImage() {
		return containerImage;
	}

	public ContainerRegistryConfiguration getRegistryConf() {
		return registryConf;
	}

	public HttpHeaders getAuthHttpHeaders() {
		return authHttpHeaders;
	}

	public RestTemplate getRestTemplate() {
		return restTemplate;
	}
}
