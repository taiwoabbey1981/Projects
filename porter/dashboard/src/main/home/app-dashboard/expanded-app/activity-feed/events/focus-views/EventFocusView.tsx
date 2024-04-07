import Loading from "components/Loading";
import Spacer from "components/porter/Spacer";
import React, { useContext, useEffect, useState } from "react";
import { Context } from "shared/Context";
import api from "shared/api";
import styled from "styled-components";
import Link from "components/porter/Link";
import BuildFailureEventFocusView from "./BuildFailureEventFocusView";
import PreDeployEventFocusView from "./PredeployEventFocusView";
import _ from "lodash";
import { PorterAppDeployEvent, PorterAppEvent } from "../types";
import DeployEventFocusView from "./DeployEventFocusView";
import { LogFilterQueryParamOpts } from "../../../logs/types";

type Props = {
    eventId: string;
    appData: any;
    filterOpts?: LogFilterQueryParamOpts;
};

const EVENT_POLL_INTERVAL = 5000; // poll every 5 seconds

const EventFocusView: React.FC<Props> = ({
    eventId,
    appData,
    filterOpts,
}) => {
    const { currentProject, currentCluster } = useContext(Context);
    const [event, setEvent] = useState<PorterAppEvent | null>(null);

    useEffect(() => {
        const getEvent = async () => {
            if (currentProject == null || currentCluster == null) {
                return;
            }
            try {
                const eventResp = await api.getPorterAppEvent(
                    "<token>",
                    {},
                    {
                        project_id: currentProject.id,
                        cluster_id: currentCluster.id,
                        event_id: eventId,
                    }
                )
                const newEvent = PorterAppEvent.toPorterAppEvent(eventResp.data.event);
                setEvent(newEvent);
                if (newEvent.metadata.end_time != null) {
                    clearInterval(intervalId);
                }
            } catch (err) {
                console.log(err);
            }
        }
        const intervalId = setInterval(getEvent, EVENT_POLL_INTERVAL);
        getEvent();
        return () => clearInterval(intervalId);
    }, []);

    const getEventFocusView = (event: PorterAppEvent, appData: any) => {
        switch (event.type) {
            case "BUILD":
                return <BuildFailureEventFocusView event={event} appData={appData} />
            case "PRE_DEPLOY":
                return <PreDeployEventFocusView event={event} appData={appData} />
            case "DEPLOY":
                return <DeployEventFocusView
                    event={event as PorterAppDeployEvent}
                    appData={appData}
                    filterOpts={filterOpts}
                />
            default:
                return null
        }
    }

    return (
        <AppearingView>
            <Link to={`/apps/${appData.app.name}/activity`}>
                <BackButton>
                    <i className="material-icons">keyboard_backspace</i>
                    Activity feed
                </BackButton>
            </Link>
            <Spacer y={0.5} />
            {event == null && <Loading />}
            {event != null && getEventFocusView(event, appData)}
        </AppearingView>
    );
};

export default EventFocusView;

export const AppearingView = styled.div`
    width: 100%;
    animation: fadeIn 0.3s 0s;
    @keyframes fadeIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
    }
`;

const BackButton = styled.div`
  display: flex;
  align-items: center;
  max-width: fit-content;
  cursor: pointer;
  font-size: 11px;
  max-height: fit-content;
  padding: 5px 13px;
  border: 1px solid #ffffff55;
  border-radius: 100px;
  color: white;
  background: #ffffff11;

  :hover {
    background: #ffffff22;
  }

  > i {
    color: white;
    font-size: 16px;
    margin-right: 6px;
  }
`;