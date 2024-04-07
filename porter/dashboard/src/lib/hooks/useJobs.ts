import _ from "lodash";
import { useEffect, useState } from "react";
import api from "shared/api";
import { useRevisionList } from "./useRevisionList";
import { z } from "zod";
import { useQuery } from "@tanstack/react-query";

const jobRunValidator = z.object({
    metadata: z.object({
        labels: z.object({
            "porter.run/app-revision-id": z.string(),
            "porter.run/service-name": z.string(),
            "porter.run/app-id": z.string(),
        }),
        creationTimestamp: z.string(),
        uid: z.string(),
    }),
    status: z.object({
        startTime: z.string().optional(),
        completionTime: z.string().optional(),
        conditions: z.array(z.object({
            lastTransitionTime: z.string(),
        })).default([]),
        succeeded: z.number().optional(),
        failed: z.number().optional(),
    }),
    revisionNumber: z.number().optional(),
    jobName: z.string().optional(),
});

export type JobRun = z.infer<typeof jobRunValidator>;

export const useJobs = (
    {
        appName,
        projectId,
        clusterId,
        deploymentTargetId,
        selectedJobName,
    }: {
        appName: string,
        projectId: number,
        clusterId: number,
        deploymentTargetId: string,
        selectedJobName: string,
    }
): {
    jobRuns: JobRun[],
    isLoadingJobRuns: boolean,
} => {
    const [jobRuns, setJobRuns] = useState<JobRun[]>([]);

    const { revisionIdToNumber } = useRevisionList({ appName, deploymentTargetId, projectId, clusterId });

    const { data, isLoading: isLoadingJobRuns } = useQuery(
        ["jobRuns", appName, deploymentTargetId, revisionIdToNumber, selectedJobName],
        async () => {
            const res = await api.appJobs(
                "<token>",
                {
                    deployment_target_id: deploymentTargetId,
                    job_name: selectedJobName === "all" ? "" : selectedJobName,
                },
                {
                    project_id: projectId,
                    cluster_id: clusterId,
                    porter_app_name: appName,
                });
            const parsed = await z.array(jobRunValidator).parseAsync(res.data);
            const parsedWithRevision = parsed.map((jobRun) => {
                const revisionId = jobRun.metadata.labels["porter.run/app-revision-id"];
                const revisionNumber = revisionIdToNumber[revisionId];
                return {
                    ...jobRun,
                    revisionNumber,
                    jobName: jobRun.metadata.labels["porter.run/service-name"],
                };
            });
            return parsedWithRevision;
        },
        {
            enabled: revisionIdToNumber != null,
            refetchInterval: 5000,
            refetchOnWindowFocus: false,
        },
    );

    useEffect(() => {
        if (data != null) {
            setJobRuns(data);
        }
    }, [data]);

    return {
        jobRuns,
        isLoadingJobRuns,
    };
};