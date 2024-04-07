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

package org.springframework.cloud.dataflow.rest.resource;

import java.util.ArrayList;
import java.util.List;

import org.springframework.cloud.dataflow.core.ConfigurationMetadataPropertyEntity;
import org.springframework.cloud.dataflow.core.Launcher;
import org.springframework.hateoas.PagedModel;
import org.springframework.hateoas.RepresentationModel;

/**
 * @author Ilayaperumal Gopinathan
 */
public class LauncherResource extends RepresentationModel<LauncherResource> {

	private String name;

	private String type;

	private String description;

	private List<ConfigurationMetadataPropertyEntity> options = new ArrayList<>();

	@SuppressWarnings("unused")
	private LauncherResource() {
	}

	public LauncherResource(Launcher launcher) {
		this.name = launcher.getName();
		this.type = launcher.getType();
		this.description = launcher.getDescription();
		this.options = launcher.getOptions();
	}

	public String getName() {
		return name;
	}

	public void setName(String name) {
		this.name = name;
	}

	public String getType() {
		return type;
	}

	public String getDescription() {
		return description;
	}

	public void setDescription(String description) {
		this.description = description;
	}

	public List<ConfigurationMetadataPropertyEntity> getOptions() {
		return options;
	}

	public void setOptions(List<ConfigurationMetadataPropertyEntity> options) {
		this.options = options;
	}

	public static class Page extends PagedModel<LauncherResource> {

	}
}
