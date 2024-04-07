/*
 * Copyright 2016-2017 the original author or authors.
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

package org.springframework.cloud.dataflow.completion;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

import org.springframework.cloud.dataflow.configuration.metadata.ApplicationConfigurationMetadataResolver;
import org.springframework.cloud.dataflow.core.AppRegistration;
import org.springframework.cloud.dataflow.core.ApplicationType;
import org.springframework.cloud.dataflow.core.TaskDefinition;
import org.springframework.cloud.dataflow.registry.service.AppRegistryService;

/**
 * Adds missing application configuration properties at the end of a well formed task
 * definition.
 *
 * @author Eric Bottard
 * @author Mark Fisher
 * @author Andy Clement
 * @author Oleg Zhurakousky
 */
class AddAppOptionsTaskExpansionStrategy implements TaskExpansionStrategy {

	private final ProposalsCollectorSupportUtils collectorSupport;

	public AddAppOptionsTaskExpansionStrategy(AppRegistryService appRegistry,
			ApplicationConfigurationMetadataResolver metadataResolver) {
		this.collectorSupport = new ProposalsCollectorSupportUtils(appRegistry, metadataResolver);
	}

	@Override
	public boolean addProposals(String text, TaskDefinition taskDefinition, int detailLevel,
			List<CompletionProposal> collector) {
		String appName = taskDefinition.getRegisteredAppName();

		AppRegistration appRegistration = this.collectorSupport.findAppRegistration(appName, ApplicationType.task);

		if (appRegistration != null) {
			Set<String> alreadyPresentOptions = new HashSet<>(taskDefinition.getProperties().keySet());
			this.collectorSupport.addPropertiesProposals(text, "", appRegistration, alreadyPresentOptions, collector, detailLevel);
		}
		return false;
	}
}
