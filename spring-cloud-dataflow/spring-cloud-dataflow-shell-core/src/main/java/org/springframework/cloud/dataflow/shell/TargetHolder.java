/*
 * Copyright 2016-2022 the original author or authors.
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

package org.springframework.cloud.dataflow.shell;

import org.springframework.util.Assert;

/**
 * A target holder, wrapping a {@link Target} that encapsulates not only the Target URI
 * but also success/error messages + status.
 *
 * @author Gunnar Hillert
 * @see Target
 * @since 1.0
 */
public class TargetHolder {

	private Target target;

	/**
	 * Constructor.
	 */
	public TargetHolder() {
	}

	/**
	 * @return the {@link Target} which encapsulates not only the Target URI but also success/error messages + status.
	 */
	public Target getTarget() {
		return target;
	}

	/**
	 * Set the Dataflow Server {@link Target}.
	 *
	 * @param target Must not be null.
	 */
	public void setTarget(Target target) {
		Assert.notNull(target, "The provided target must not be null.");
		this.target = target;
	}

}
