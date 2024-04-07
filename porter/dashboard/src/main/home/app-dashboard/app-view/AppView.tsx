import React, { useMemo } from "react";
import { withRouter, type RouteComponentProps } from "react-router";
import styled from "styled-components";
import { z } from "zod";

import Back from "components/porter/Back";

import AppDataContainer from "./AppDataContainer";
import AppHeader from "./AppHeader";
import { LatestRevisionProvider } from "./LatestRevisionContext";

export const porterAppValidator = z.object({
  id: z.number(),
  name: z.string(),
  git_branch: z.string().optional(),
  git_repo_id: z.number().optional(),
  repo_name: z.string().optional(),
  build_context: z.string().optional(),
  builder: z.string().optional(),
  buildpacks: z
    .string()
    .transform((s) => s.split(","))
    .optional(),
  dockerfile: z.string().optional(),
  image_repo_uri: z.string().optional(),
  porter_yaml_path: z.string().optional(),
  pull_request_url: z.string().optional(),
});
export type PorterAppRecord = z.infer<typeof porterAppValidator>;

type Props = RouteComponentProps;

const AppView: React.FC<Props> = ({ match }) => {
  const params = useMemo(() => {
    const { params } = match;
    const validParams = z
      .object({
        appName: z.string(),
        tab: z.string().optional(),
      })
      .safeParse(params);

    if (!validParams.success) {
      return {
        appName: undefined,
        tab: undefined,
      };
    }

    return validParams.data;
  }, [match]);

  return (
    <LatestRevisionProvider appName={params.appName}>
      <StyledExpandedApp>
        <Back to="/apps" />
        <AppHeader />
        <AppDataContainer tabParam={params.tab} />
      </StyledExpandedApp>
    </LatestRevisionProvider>
  );
};

export default withRouter(AppView);

const StyledExpandedApp = styled.div`
  width: 100%;
  height: 100%;

  animation: fadeIn 0.5s 0s;
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
`;
