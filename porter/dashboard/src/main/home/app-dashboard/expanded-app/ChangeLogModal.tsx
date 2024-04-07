import React, { useContext, useEffect, useRef, useState } from "react";
import styled from "styled-components";
import Modal from "components/porter/Modal";
import Loading from "components/Loading";
import Text from "components/porter/Text";
import yaml from "js-yaml";
import DiffViewer, { DiffMethod } from "react-diff-viewer";
import Button from "components/porter/Button";
import ConfirmOverlay from "components/porter/ConfirmOverlay";
import Spacer from "components/porter/Spacer";
import Checkbox from "components/porter/Checkbox";
import { ChartType } from "shared/types";
import * as Diff from "deep-diff";
import api from "shared/api";
import { Context } from "shared/Context";
import ChangeLogComponent from "./ChangeLogComponent";

type Props = {
  modalVisible: boolean;
  setModalVisible: (x: boolean) => void;
  revision: number;
  currentChart: ChartType;
  revertModal?: boolean;
  appData: any;
  diffContent: boolean;
  setDiffContent: (x: boolean) => void;
};

const ChangeLogModal: React.FC<Props> = ({
  revision,
  appData,
  currentChart,
  modalVisible,
  revertModal,
  setModalVisible,
}) => {
  const [values, setValues] = useState("");
  const [chartEvent, setChartEvent] = useState(null);
  const [eventValues, setEventValues] = useState("");
  const [prevChartEvent, setPrevChartEvent] = useState(null);
  const [prevEventValues, setPrevEventValues] = useState("");
  const [showRawDiff, setShowRawDiff] = useState(false);
  const [showOverlay, setShowOverlay] = useState<boolean>(false);
  const [changesConfig, setChangesConfig] = useState<boolean>(true);

  const [loading, setLoading] = useState(false);
  const { currentCluster, currentProject, setCurrentError } = useContext(
    Context
  );
  useEffect(() => {
    let values = "# Nothing here yet";
    if (currentChart.config) {
      values = yaml.dump(currentChart.config);
    }
    setValues(values);
  }, [currentChart.config]); // It will run this effect whenever currentChart.config changes

  const getChartData = async (chart: ChartType) => {
    setLoading(true);
    const res = await api.getChart(
      "<token>",
      {},
      {
        name: chart.name,
        namespace: chart.namespace,
        cluster_id: currentCluster.id,
        revision: revision,
        id: currentProject.id,
      }
    );
    const updatedChart = res.data;
    setLoading(false);
    return updatedChart;
  };

  const revertToRevision = async (revision: number) => {
    setLoading(true);
    try {
      await api
        .rollbackPorterApp(
          "<token>",
          {
            revision,
          },
          {
            project_id: appData.app.project_id,
            stack_name: appData.app.name,
            cluster_id: appData.app.cluster_id,
          }
        )
      window.location.reload();
    } catch (err) {
      setLoading(false);
      console.log(err)
    }
  }

  const getPrevChartData = async (chart: ChartType) => {
    setLoading(true);
    const prevRevision = revision - 1;
    const res = await api.getChart(
      "<token>",
      {},
      {
        name: chart.name,
        namespace: chart.namespace,
        cluster_id: currentCluster.id,
        revision: prevRevision,
        id: currentProject.id,
      }
    );
    const updatedChart = res.data;
    setLoading(false);
    return updatedChart;
  };

  useEffect(() => {
    const fetchData = async () => {
      // Fetch the chart data
      const updatedChart = await getChartData(currentChart);
      const prevChart = await getPrevChartData(currentChart);

      // Now that we've waited for getChartData to finish, process the result
      let eventValues = "# Nothing here yet";
      if (updatedChart?.config) {
        eventValues = yaml.dump(updatedChart?.config);
      }
      let prevEventValues = "# Nothing here yet";
      if (prevChart?.config) {
        prevEventValues = yaml.dump(prevChart?.config);
      }
      setEventValues(eventValues);
      setChartEvent(updatedChart);
      setPrevEventValues(prevEventValues);
      setPrevChartEvent(prevChart);
    };

    fetchData();
  }, [currentChart.config]);

  return (
    <>
      <Modal closeModal={() => setModalVisible(false)} width={"800px"}>
        {revertModal ? <Text size={18}> Revert to version no. {revision} </Text> : <Text size={18}>Changes for version no. {revision}</Text>}
        <Spacer y={1} />
        {loading ? (
          <Loading /> // <-- Render loading state
        ) : (
          revertModal ? (<>
            <div style={{ maxHeight: "400px", overflowY: "auto", borderRadius: "8px" }}>
              <DiffViewer
                leftTitle={revertModal ? `Current Revision` : `Revision No. ${revision - 1}`}
                rightTitle={`Revision No. ${revision}`}
                oldValue={revertModal ? values : prevEventValues}
                newValue={eventValues}
                splitView={true}
                hideLineNumbers={false}
                useDarkTheme={true}
                compareMethod={DiffMethod.TRIMMED_LINES}
              />
            </div>
          </>) :
            (<>
              {showRawDiff ? (
                <>
                  <div style={{ maxHeight: "400px", overflowY: "auto", borderRadius: "8px" }}>
                    <DiffViewer
                      leftTitle={revertModal ? `Current Revision` : `Revision No. ${revision - 1}`}
                      rightTitle={`Revision No. ${revision}`}
                      oldValue={revertModal ? values : prevEventValues}
                      newValue={eventValues}
                      splitView={true}
                      hideLineNumbers={false}
                      useDarkTheme={true}
                      compareMethod={DiffMethod.TRIMMED_LINES}
                    />
                  </div>
                </>
              ) : (
                <div style={{ maxHeight: "400px", overflowY: "auto" }}>
                  {revertModal ?

                    <ChangeLogComponent
                      oldYaml={currentChart.config}
                      newYaml={chartEvent?.config}
                      appData={appData}
                    />
                    : <ChangeLogComponent
                      oldYaml={prevChartEvent?.config}
                      newYaml={chartEvent?.config}
                      appData={appData}
                    />

                  }
                </div>
              )}

              {changesConfig && (<>
                <Spacer y={1} />
                <div style={{ display: "flex" }}>

                  <Checkbox
                    checked={showRawDiff}
                    toggleChecked={() => setShowRawDiff(!showRawDiff)}
                  >
                    <Text>Show raw diff</Text>
                  </Checkbox>
                </div></>)}
            </>))
        }

        {revertModal && (
          <>
            <Spacer y={1} />
            <Button
              onClick={() => setShowOverlay(true)}
              width={"110px"}
              loadingText={"Submitting..."}
            >
              Revert
            </Button>
          </>
        )}
        {showOverlay && (

          <ConfirmOverlay
            loading={loading}
            message={`Are you sure you want to revert to version no. ${revision}?`}
            onYes={() => revertToRevision(revision)}
            onNo={() => setShowOverlay(false)}
          />

        )}

      </Modal>
    </>
  );
};

export default ChangeLogModal;