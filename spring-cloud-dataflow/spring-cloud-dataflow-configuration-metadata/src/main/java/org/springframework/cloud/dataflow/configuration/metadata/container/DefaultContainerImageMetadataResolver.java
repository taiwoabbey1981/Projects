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

package org.springframework.cloud.dataflow.configuration.metadata.container;

import java.util.Collections;
import java.util.Map;

import org.springframework.cloud.dataflow.container.registry.ContainerRegistryException;
import org.springframework.cloud.dataflow.container.registry.ContainerRegistryRequest;
import org.springframework.cloud.dataflow.container.registry.ContainerRegistryService;
import org.springframework.util.StringUtils;

/**
 * Leverages the Docker Registry HTTP V2 API to retrieve the configuration object and the labels
 * form the specified image.
 *
 * @author Christian Tzolov
 * @author Ilayaperumal Gopinathan
 */
public class DefaultContainerImageMetadataResolver implements ContainerImageMetadataResolver {

	private final ContainerRegistryService containerRegistryService;

	public DefaultContainerImageMetadataResolver(ContainerRegistryService containerRegistryService) {
		this.containerRegistryService = containerRegistryService;
	}

	@Override
	public Map<String, String> getImageLabels(String imageName) {

		if (!StringUtils.hasText(imageName)) {
			throw new ContainerRegistryException("Null or empty image name");
		}

		ContainerRegistryRequest registryRequest = this.containerRegistryService.getRegistryRequest(imageName);

		Map manifest = this.containerRegistryService.getImageManifest(registryRequest, Map.class);

		if (manifest != null && !isNotNullMap(manifest.get("config"))) {
			throw new ContainerRegistryException(
					String.format("Image [%s] has incorrect or missing manifest config element: %s",
							imageName, manifest.toString()));
		}
		if (manifest != null) {
			String configDigest = ((Map<String, String>) manifest.get("config")).get("digest");

			if (!StringUtils.hasText(configDigest)) {
				throw new ContainerRegistryException(
						String.format("Missing or invalid Configuration Digest: [%s] for image [%s]", configDigest,
								imageName));
			}

			Map<String, Object> configBlobMap = this.containerRegistryService.getImageBlob(registryRequest, configDigest, Map.class);

			if (configBlobMap == null) {
				throw new ContainerRegistryException(
						String.format("Failed to retrieve configuration json for image [%s] with digest [%s]",
								imageName, configDigest));
			}

			if (!isNotNullMap(configBlobMap.get("config"))) {
				throw new ContainerRegistryException(
						String.format(
								"Configuration json for image [%s] with digest [%s] has incorrect Config Blog element",
								imageName, configDigest));
			}

			Map<String, Object> configElement = (Map<String, Object>) configBlobMap.get("config");

			return isNotNullMap(configElement.get("Labels")) ?
					(Map<String, String>) configElement.get("Labels") : Collections.emptyMap();
		}
		else {
			throw new ContainerRegistryException(
					String.format("Image [%s] is missing manifest", imageName));
		}
	}

	private boolean isNotNullMap(Object object) {
		return object != null && (object instanceof Map);
	}
}
