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

package org.springframework.cloud.dataflow.core;

import java.util.ArrayList;
import java.util.List;

import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

/**
 * @author Christian Tzolov
 * @author Ilayaperumal Gopinathan
 */
public class ArgumentSanitizerTest {

	private ArgumentSanitizer sanitizer;

	private static final String[] keys = { "password", "secret", "key", "token", ".*credentials.*",
			"vcap_services", "url" };

	@Before
	public void before() {
		sanitizer = new ArgumentSanitizer();
	}

	@Test
	public void testSanitizeProperties() {
		for (String key : keys) {
			Assert.assertEquals("--" + key + "=******", sanitizer.sanitize("--" + key + "=foo"));
			Assert.assertEquals("******", sanitizer.sanitize(key, "bar"));
		}
	}

	@Test
	public void testSanitizeArguments() {
		final List<String> arguments = new ArrayList<>();

		for (String key : keys) {
			arguments.add("--" + key + "=foo");
		}

		final List<String> sanitizedArguments = sanitizer.sanitizeArguments(arguments);

		Assert.assertEquals(keys.length, sanitizedArguments.size());

		int order = 0;
		for(String sanitizedString : sanitizedArguments) {
			Assert.assertEquals("--" + keys[order] + "=******", sanitizedString);
			order++;
		}
	}


	@Test
	public void testMultipartProperty() {
		Assert.assertEquals("--password=******", sanitizer.sanitize("--password=boza"));
		Assert.assertEquals("--one.two.password=******", sanitizer.sanitize("--one.two.password=boza"));
		Assert.assertEquals("--one_two_password=******", sanitizer.sanitize("--one_two_password=boza"));
	}

//	@Test
//	public void testHierarchicalPropertyNames() {
//		Assert.assertEquals("time --password='******' | log",
//				sanitizer.(new StreamDefinition("stream", "time --password=bar | log")));
//	}
//
//	@Test
//	public void testStreamPropertyOrder() {
//		Assert.assertEquals("time --some.password='******' --another-secret='******' | log",
//				sanitizer.sanitizeStream(new StreamDefinition("stream", "time --some.password=foobar --another-secret=kenny | log")));
//	}
//
//	@Test
//	public void testStreamMatcherWithHyphenDotChar() {
//		Assert.assertEquals("twitterstream --twitter.credentials.access-token-secret='******' "
//						+ "--twitter.credentials.access-token='******' --twitter.credentials.consumer-secret='******' "
//						+ "--twitter.credentials.consumer-key='******' | "
//						+ "filter --expression=#jsonPath(payload,'$.lang')=='en' | "
//						+ "twitter-sentiment --vocabulary=https://dl.bintray.com/test --model-fetch=output/test "
//						+ "--model=https://dl.bintray.com/test | field-value-counter --field-name=sentiment --name=sentiment",
//				sanitizer.sanitizeStream(new StreamDefinition("stream", "twitterstream "
//						+ "--twitter.credentials.consumer-key=dadadfaf --twitter.credentials.consumer-secret=dadfdasfdads "
//						+ "--twitter.credentials.access-token=58849055-dfdae "
//						+ "--twitter.credentials.access-token-secret=deteegdssa4466 | filter --expression='#jsonPath(payload,''$.lang'')==''en''' | "
//						+ "twitter-sentiment --vocabulary=https://dl.bintray.com/test --model-fetch=output/test --model=https://dl.bintray.com/test | "
//						+ "field-value-counter --field-name=sentiment --name=sentiment")));
//	}
//
//	@Test
//	public void testStreamSanitizeOriginalDsl() {
//		StreamDefinition streamDefinition = new StreamDefinition("test", "time --password='******' | log --password='******'", "time --password='******' | log");
//		Assert.assertEquals("time --password='******' | log", sanitizer.sanitizeOriginalStreamDsl(streamDefinition));
//	}
}
