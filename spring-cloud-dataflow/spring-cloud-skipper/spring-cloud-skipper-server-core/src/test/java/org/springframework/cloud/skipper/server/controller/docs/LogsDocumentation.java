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

package org.springframework.cloud.skipper.server.controller.docs;

import java.nio.charset.Charset;
import java.util.Collections;

import org.junit.Test;

import org.springframework.cloud.skipper.domain.LogInfo;
import org.springframework.cloud.skipper.domain.Release;
import org.springframework.http.MediaType;

import static org.mockito.Mockito.when;
import static org.springframework.restdocs.mockmvc.RestDocumentationRequestBuilders.get;
import static org.springframework.restdocs.payload.PayloadDocumentation.responseBody;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

/**
 * @author Ilayaperumal Gopinathan
 */
public class LogsDocumentation extends BaseDocumentation {

	@Test
	public void getLogsofRelease() throws Exception {
		Release release = createTestRelease();
		when(this.releaseService.getLog(release.getName())).thenReturn(new LogInfo(Collections.emptyMap()));
		final MediaType contentType = new MediaType(MediaType.APPLICATION_JSON.getType(),
				MediaType.APPLICATION_JSON.getSubtype(), Charset.forName("utf8"));

		this.mockMvc.perform(
				get("/api/release/logs/{releaseName}", release.getName()).accept(MediaType.APPLICATION_JSON)
						.contentType(contentType))
				.andDo(print())
				.andExpect(status().isOk())
				.andDo(this.documentationHandler.document(
						responseBody()));
	}

	@Test
	public void getLogsofReleaseByAppName() throws Exception {
		Release release = createTestRelease();
		when(this.releaseService.getLog(release.getName(), "myapp")).thenReturn(new LogInfo(Collections.EMPTY_MAP));

		this.mockMvc.perform(
				get("/api/release/logs/{releaseName}/{appName}",
						release.getName(), "myapp"))
				.andDo(print())
				.andExpect(status().isOk())
				.andDo(this.documentationHandler.document(
						responseBody()));
	}
}
