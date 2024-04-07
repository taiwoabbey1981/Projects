import React, { useState } from "react";
import styled from "styled-components";

import build from "assets/build.png";

import run_for from "assets/run_for.png";
import refresh from "assets/refresh.png";

import Text from "components/porter/Text";
import Container from "components/porter/Container";
import Spacer from "components/porter/Spacer";
import Link from "components/porter/Link";
import Icon from "components/porter/Icon";
import api from "shared/api";
import { Log } from "main/home/cluster-dashboard/expanded-chart/logs-section/useAgentLogs";
import JSZip from "jszip";
import Anser, { AnserJsonEntry } from "anser";
import { getDuration, getStatusColor, getStatusIcon, triggerWorkflow } from '../utils';
import { StyledEventCard } from "./EventCard";
import document from "assets/document.svg";
import { PorterAppEvent } from "../types";

type Props = {
  event: PorterAppEvent;
  appData: any;
};

const BuildEventCard: React.FC<Props> = ({ event, appData }) => {
  const renderStatusText = (event: PorterAppEvent) => {
    switch (event.status) {
      case "SUCCESS":
        return <Text color={getStatusColor(event.status)}>Build succeeded</Text>;
      case "FAILED":
        return <Text color={getStatusColor(event.status)}>Build failed</Text>;
      default:
        return <Text color={getStatusColor(event.status)}>Build in progress...</Text>;
    }
  };

  const renderInfoCta = (event: PorterAppEvent) => {
    switch (event.status) {
      case "SUCCESS":
        return null;
      case "FAILED":
        return (
          <Wrapper>
            <Link to={`/apps/${appData.app.name}/events?event_id=${event.id}`} hasunderline>
              <Container row>
                <Icon src={document} height="10px" />
                <Spacer inline width="5px" />
                View details
              </Container>
            </Link>
            <Spacer inline x={1} />
            <Link hasunderline onClick={() => triggerWorkflow(appData)}>
              <Container row>
                <Icon height="10px" src={refresh} />
                <Spacer inline width="5px" />
                Retry
              </Container>
            </Link>
          </Wrapper>
        );
      default:
        return (
          <Wrapper>
            <Link
              hasunderline
              target="_blank"
              to={`https://github.com/${appData.app.repo_name}/actions/runs/${event.metadata?.action_run_id}`}
            >
              View live logs
            </Link>
            <Spacer inline x={1} />
          </Wrapper>
        );
    }
  };

  return (
    <StyledEventCard>
      <Container row spaced>
        <Container row>
          <Icon height="16px" src={build} />
          <Spacer inline width="10px" />
          <Text>Application build</Text>
        </Container>
        <Container row>
          <Icon height="14px" src={run_for} />
          <Spacer inline width="6px" />
          <Text color="helper">{getDuration(event)}</Text>
        </Container>
      </Container>
      <Spacer y={0.5} />
      <Container row spaced>
        <Container row>
          <Icon height="12px" src={getStatusIcon(event.status)} />
          <Spacer inline width="10px" />
          {renderStatusText(event)}
          <Spacer inline x={1} />
          {renderInfoCta(event)}
          <Spacer inline x={1} />
        </Container>
      </Container>
    </StyledEventCard>
  );
};

export default BuildEventCard;

const Wrapper = styled.div`
  margin-top: -3px;
`;
