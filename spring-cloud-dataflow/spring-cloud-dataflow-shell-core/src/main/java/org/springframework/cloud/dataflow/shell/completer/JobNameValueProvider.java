/*
 * Copyright 2022 the original author or authors.
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

package org.springframework.cloud.dataflow.shell.completer;

import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

import org.springframework.cloud.dataflow.rest.resource.JobExecutionResource;
import org.springframework.cloud.dataflow.shell.config.DataFlowShell;
import org.springframework.shell.CompletionContext;
import org.springframework.shell.CompletionProposal;
import org.springframework.shell.standard.ValueProvider;
import org.springframework.stereotype.Component;

/**
 * A {@link org.springframework.shell.standard.ValueProvider} that provides completion for Data Flow job execution names.
 *
 * @author Chris Bono
 */
@Component
public class JobNameValueProvider implements ValueProvider {

	private final DataFlowShell dataFlowShell;

	public JobNameValueProvider(DataFlowShell dataFlowShell) {
		this.dataFlowShell = dataFlowShell;
	}

	@Override
	public List<CompletionProposal> complete(CompletionContext completionContext) {
		return StreamSupport.stream(
				dataFlowShell.getDataFlowOperations().jobOperations().executionList().spliterator(), false)
				.map(JobExecutionResource::getName)
				.filter(jobExecName -> !"?".equals(jobExecName))
				.map(CompletionProposal::new)
				.collect(Collectors.toList());
	}
}
