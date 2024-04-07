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
package org.springframework.cloud.dataflow.audit.service;

import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.springframework.cloud.dataflow.audit.repository.AuditRecordRepository;
import org.springframework.cloud.dataflow.core.ArgumentSanitizer;
import org.springframework.cloud.dataflow.core.AuditActionType;
import org.springframework.cloud.dataflow.core.AuditOperationType;
import org.springframework.cloud.dataflow.core.AuditRecord;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.util.Assert;

/**
 * Default implementation of the {@link AuditRecordService}.
 *
 * @author Gunnar Hillert
 * @author Daniel Serleg
 * @author Corneil du Plessis
 */
public class DefaultAuditRecordService implements AuditRecordService {

    private static final Logger logger = LoggerFactory.getLogger(DefaultAuditRecordService.class);

    private final AuditRecordRepository auditRecordRepository;

    private final ObjectMapper objectMapper;

    private final ArgumentSanitizer sanitizer;

    public DefaultAuditRecordService(AuditRecordRepository auditRecordRepository) {

        this(auditRecordRepository, null);
    }

    public DefaultAuditRecordService(AuditRecordRepository auditRecordRepository, ObjectMapper objectMapper) {

        Assert.notNull(auditRecordRepository, "auditRecordRepository must not be null.");
        this.auditRecordRepository = auditRecordRepository;
        if (objectMapper == null) {
            objectMapper = new ObjectMapper();
            objectMapper.configure(SerializationFeature.FAIL_ON_EMPTY_BEANS, false);
        }
        this.objectMapper = objectMapper;
        this.sanitizer = new ArgumentSanitizer();
    }

    @Override
    public AuditRecord populateAndSaveAuditRecord(AuditOperationType auditOperationType,
                                                  AuditActionType auditActionType,
                                                  String correlationId, String data, String platformName) {

        Assert.notNull(auditActionType, "auditActionType must not be null.");
        Assert.notNull(auditOperationType, "auditOperationType must not be null.");

        final AuditRecord auditRecord = new AuditRecord();
        auditRecord.setAuditAction(auditActionType);
        auditRecord.setAuditOperation(auditOperationType);
        auditRecord.setCorrelationId(correlationId);
        auditRecord.setAuditData(data);
        auditRecord.setPlatformName(platformName);
        return this.auditRecordRepository.save(auditRecord);
    }

    @Override
    public AuditRecord populateAndSaveAuditRecordUsingMapData(
            AuditOperationType auditOperationType,
            AuditActionType auditActionType,
            String correlationId, Map<String, Object> data,
            String platformName
    ) {

        String dataAsString;
        try {
            Map<String, Object> sanitizedData = sanitizeMap(data);
            dataAsString = objectMapper.writeValueAsString(sanitizedData);
        } catch (JsonProcessingException e) {
            logger.error("Error serializing audit record data.  Data = " + data);
            dataAsString = "Error serializing audit record data.  Data = " + data;
        }
        return this.populateAndSaveAuditRecord(auditOperationType, auditActionType, correlationId, dataAsString, platformName);
    }

    private Map<String, Object> sanitizeMap(Map<String, Object> data) {

        final Map<String, Object> result = new HashMap<>();
        data.forEach((k, v) -> result.put(k, sanitize(k, v)));
        return result;
    }

    private Object sanitize(String key, Object value) {

        if (value instanceof String) {
            return sanitizer.sanitize(key, (String) value);
        } else if (value instanceof Map) {
            Map<String, Object> input = (Map) value;
            return sanitizeMap(input);
        } else {
            return value;
        }
    }


    @Override
    public Page<AuditRecord> findAuditRecordByAuditOperationTypeAndAuditActionTypeAndDate(
            Pageable pageable,
            AuditActionType[] actions,
            AuditOperationType[] operations,
            Instant fromDate,
            Instant toDate) {

        return this.auditRecordRepository.findByActionTypeAndOperationTypeAndDate(operations, actions, fromDate, toDate,
                pageable);
    }

    @Override
    public Optional<AuditRecord> findById(Long id) {

        return this.auditRecordRepository.findById(id);
    }

}
