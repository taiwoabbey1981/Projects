/*
 * Copyright 2016-2018 the original author or authors.
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

package org.springframework.cloud.dataflow.server.controller;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.Map;
import java.util.Optional;

import org.json.JSONObject;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.junit.runner.RunWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import org.springframework.cloud.dataflow.core.ApplicationType;
import org.springframework.cloud.dataflow.core.StreamAppDefinition;
import org.springframework.cloud.dataflow.core.StreamDefinition;
import org.springframework.cloud.dataflow.core.StreamDefinitionService;
import org.springframework.cloud.dataflow.core.StreamDeployment;
import org.springframework.cloud.dataflow.rest.SkipperStream;
import org.springframework.cloud.dataflow.rest.UpdateStreamRequest;
import org.springframework.cloud.dataflow.rest.resource.StreamDeploymentResource;
import org.springframework.cloud.dataflow.server.repository.StreamDefinitionRepository;
import org.springframework.cloud.dataflow.server.service.StreamService;
import org.springframework.cloud.deployer.spi.app.DeploymentState;
import org.springframework.cloud.skipper.domain.Deployer;
import org.springframework.cloud.skipper.domain.PackageIdentifier;
import org.springframework.mock.web.MockHttpServletRequest;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyList;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.ArgumentMatchers.isA;
import static org.mockito.ArgumentMatchers.isNull;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

/**
 * Unit tests for SkipperStreamDeploymentController.
 *
 * @author Eric Bottard
 * @author Ilayaperumal Gopinathan
 * @author Christian Tzolov
 */
@RunWith(MockitoJUnitRunner.class)
public class StreamDeploymentControllerTests {

	@Rule
	public ExpectedException thrown = ExpectedException.none();

	private StreamDeploymentController controller;

	@Mock
	private StreamDefinitionRepository streamDefinitionRepository;

	@Mock
	private StreamService streamService;

	@Mock
	private StreamDefinitionService streamDefinitionService;

	@Mock
	private Deployer deployer;

	@Before
	public void setup() {
		MockHttpServletRequest request = new MockHttpServletRequest();
		RequestContextHolder.setRequestAttributes(new ServletRequestAttributes(request));
		this.controller = new StreamDeploymentController(streamDefinitionRepository, streamService, streamDefinitionService);
	}

	@Test
	@SuppressWarnings("unchecked")
	public void testDeployViaStreamService() {
		this.controller.deploy("test", new HashMap<>());
		ArgumentCaptor<String> argumentCaptor1 = ArgumentCaptor.forClass(String.class);
		ArgumentCaptor<Map> argumentCaptor2 = ArgumentCaptor.forClass(Map.class);
		verify(streamService).deployStream(argumentCaptor1.capture(), argumentCaptor2.capture());
		Assert.assertEquals(argumentCaptor1.getValue(), "test");
	}

	@Test
	public void testScaleApplicationInstances() {
		this.controller.scaleApplicationInstances("ticktock", "time", 666, null);
		verify(streamService).scaleApplicationInstances(eq("ticktock"), eq("time"), eq(666), isNull());

		this.controller.scaleApplicationInstances("stream", "foo", 2, new HashMap<>());
		verify(streamService).scaleApplicationInstances(eq("stream"), eq("foo"), eq(2), isA(Map.class));

		this.controller.scaleApplicationInstances("stream", "bar", 3,
				Collections.singletonMap("key", "value"));
		verify(streamService).scaleApplicationInstances(eq("stream"), eq("bar"), eq(3),
				eq(Collections.singletonMap("key", "value")));
	}

	@Test
	public void testUpdateStream() {
		Map<String, String> deploymentProperties = new HashMap<>();
		deploymentProperties.put(SkipperStream.SKIPPER_PACKAGE_NAME, "ticktock");
		deploymentProperties.put(SkipperStream.SKIPPER_PACKAGE_VERSION, "1.0.0");
		deploymentProperties.put("version.log", "1.2.0.RELEASE");

		UpdateStreamRequest updateStreamRequest = new UpdateStreamRequest("ticktock", new PackageIdentifier(), deploymentProperties);
		this.controller.update("ticktock", updateStreamRequest);
		ArgumentCaptor<UpdateStreamRequest> argumentCaptor1 = ArgumentCaptor.forClass(UpdateStreamRequest.class);
		verify(streamService).updateStream(eq("ticktock"), argumentCaptor1.capture());
		Assert.assertEquals(updateStreamRequest, argumentCaptor1.getValue());
	}

	@Test
	public void testStreamManifest() {
		this.controller.manifest("ticktock", 666);
		verify(streamService, times(1)).manifest(eq("ticktock"), eq(666));
	}

	@Test
	public void testStreamHistory() {
		this.controller.history("releaseName");
		verify(streamService, times(1)).history(eq("releaseName"));
	}

	@Test
	public void testRollbackViaStreamService() {
		this.controller.rollback("test1", 2);
		ArgumentCaptor<String> argumentCaptor1 = ArgumentCaptor.forClass(String.class);
		ArgumentCaptor<Integer> argumentCaptor2 = ArgumentCaptor.forClass(Integer.class);
		verify(streamService).rollbackStream(argumentCaptor1.capture(), argumentCaptor2.capture());
		Assert.assertEquals(argumentCaptor1.getValue(), "test1");
		Assert.assertEquals("Rollback version is incorrect", 2, (int) argumentCaptor2.getValue());
	}

	@Test
	public void testPlatformsListViaSkipperClient() {
		when(streamService.platformList()).thenReturn(Arrays.asList(deployer));
		this.controller.platformList();
		verify(streamService, times(1)).platformList();
	}

	@Test
	public void testShowStreamInfo() {
		Map<String, String> deploymentProperties1 = new HashMap<>();
		deploymentProperties1.put("test1", "value1");
		Map<String, String> deploymentProperties2 = new HashMap<>();
		deploymentProperties2.put("test2", "value2");
		Map<String, Map<String, String>> streamDeploymentProperties = new HashMap<>();
		streamDeploymentProperties.put("time", deploymentProperties1);
		streamDeploymentProperties.put("log", deploymentProperties2);
		Map<String, String> appVersions = new HashMap<>();
		appVersions.put("time", "3.2.0");
		appVersions.put("log", "3.2.0");
		StreamDefinition streamDefinition = new StreamDefinition("testStream1", "time | log");
		StreamDeployment streamDeployment = new StreamDeployment(streamDefinition.getName(),
				new JSONObject(streamDeploymentProperties).toString());
		Map<StreamDefinition, DeploymentState> streamDeploymentStates = new HashMap<>();
		streamDeploymentStates.put(streamDefinition, DeploymentState.deployed);

		StreamAppDefinition streamAppDefinition1 = new StreamAppDefinition("time", "time", ApplicationType.source, streamDefinition.getName(), new HashMap<>());
		StreamAppDefinition streamAppDefinition2 = new StreamAppDefinition("log", "log", ApplicationType.sink, streamDefinition.getName(), new HashMap<>());

		when(this.streamDefinitionRepository.findById(streamDefinition.getName())).thenReturn(Optional.of(streamDefinition));
		when(this.streamService.info(streamDefinition.getName())).thenReturn(streamDeployment);
		when(this.streamService.state(anyList())).thenReturn(streamDeploymentStates);
		LinkedList<StreamAppDefinition> streamAppDefinitions = new LinkedList<>();
		streamAppDefinitions.add(streamAppDefinition1);
		streamAppDefinitions.add(streamAppDefinition2);
		when(this.streamDefinitionService.redactDsl(any())).thenReturn("time | log");

		StreamDeploymentResource streamDeploymentResource = this.controller.info(streamDefinition.getName(), false);
		Assert.assertEquals(streamDeploymentResource.getStreamName(), streamDefinition.getName());
		Assert.assertEquals(streamDeploymentResource.getDslText(), streamDefinition.getDslText());
		Assert.assertEquals(streamDeploymentResource.getStreamName(), streamDefinition.getName());
		Assert.assertEquals("{\"log\":{\"test2\":\"value2\"},\"time\":{\"test1\":\"value1\"}}", streamDeploymentResource.getDeploymentProperties());
		Assert.assertEquals(streamDeploymentResource.getStatus(), DeploymentState.deployed.name());
	}

	@Test
	public void testReuseDeploymentProperties() {
		Map<String, String> deploymentProperties1 = new HashMap<>();
		deploymentProperties1.put("test1", "value1");
		Map<String, String> deploymentProperties2 = new HashMap<>();
		deploymentProperties2.put("test2", "value2");
		Map<String, Map<String, String>> streamDeploymentProperties = new HashMap<>();
		streamDeploymentProperties.put("time", deploymentProperties1);
		streamDeploymentProperties.put("log", deploymentProperties2);
		StreamDefinition streamDefinition = new StreamDefinition("testStream1", "time | log");
		StreamDeployment streamDeployment = new StreamDeployment(streamDefinition.getName(),
				new JSONObject(streamDeploymentProperties).toString());
		Map<StreamDefinition, DeploymentState> streamDeploymentStates = new HashMap<>();
		streamDeploymentStates.put(streamDefinition, DeploymentState.undeployed);

		when(this.streamDefinitionRepository.findById(streamDefinition.getName())).thenReturn(Optional.of(streamDefinition));
		when(this.streamService.info(streamDefinition.getName())).thenReturn(streamDeployment);
		when(this.streamService.state(anyList())).thenReturn(streamDeploymentStates);
		when(this.streamDefinitionService.redactDsl(any())).thenReturn("time | log");

		StreamDeploymentResource streamDeploymentResource = this.controller.info(streamDefinition.getName(), true);
		Assert.assertEquals(streamDeploymentResource.getStreamName(), streamDefinition.getName());
		Assert.assertEquals(streamDeploymentResource.getDslText(), streamDefinition.getDslText());
		Assert.assertEquals(streamDeploymentResource.getStreamName(), streamDefinition.getName());
		Assert.assertEquals("{\"log\":{\"test2\":\"value2\"},\"time\":{\"test1\":\"value1\"}}", streamDeploymentResource.getDeploymentProperties());
		Assert.assertEquals(streamDeploymentResource.getStatus(), DeploymentState.undeployed.name());
	}

}
