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

package org.springframework.cloud.dataflow.server.service;

import org.springframework.cloud.dataflow.core.TaskDefinition;

/**
 * Provides task saving service.
 *
 * @author Daniel Serleg
 */
public interface TaskSaveService {
    /**
     * Saves the task definition. If it is a Composed Task then the task definitions required
     * for a ComposedTaskRunner task are also created.
     *
     * @param taskDefinition the task definition object
     */
    void saveTaskDefinition(TaskDefinition taskDefinition);
}
