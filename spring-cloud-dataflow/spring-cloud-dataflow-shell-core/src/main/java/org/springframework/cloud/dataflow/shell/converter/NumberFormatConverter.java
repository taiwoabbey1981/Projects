/*
 * Copyright 2015-2022 the original author or authors.
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

package org.springframework.cloud.dataflow.shell.converter;

import java.text.DecimalFormat;
import java.text.NumberFormat;

import org.springframework.core.convert.converter.Converter;
import org.springframework.stereotype.Component;

/**
 * Knows how to convert from a String to a new {@link NumberFormat} instance.
 *
 * @author Eric Bottard
 * @author Chris Bono
 */
@Component
public class NumberFormatConverter implements Converter<String, NumberFormat> {

	public static final String DEFAULT = "<use platform locale>";

	@Override
	public NumberFormat convert(String source) {
		if (DEFAULT.equals(source)) {
			return NumberFormat.getInstance();
		}
		else {
			return new DecimalFormat(source);
		}
	}
}
