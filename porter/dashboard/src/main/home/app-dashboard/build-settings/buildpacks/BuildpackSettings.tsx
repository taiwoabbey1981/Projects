import Helper from "components/form-components/Helper";
import React, { useContext, useEffect, useMemo, useState } from "react";
import api from "shared/api";
import { Context } from "shared/Context";
import styled, { keyframes } from "styled-components";
import Button from "components/porter/Button";
import Spacer from "components/porter/Spacer";
import Error from "components/porter/Error";
import { PorterApp } from "../../types/porterApp";
import BuildpackList from "./BuildpackList";
import {
  BUILDPACK_TO_NAME,
  BuildConfig,
  Buildpack,
  DEFAULT_BUILDER_NAME,
  DEFAULT_HEROKU_STACK,
  DetectedBuildpack
} from "../../types/buildpack";
import BuildpackConfigurationModal from "./BuildpackConfigurationModal";

const BuildpackSettings: React.FC<{
  porterApp: PorterApp;
  updatePorterApp: (attrs: Partial<PorterApp>) => void;
  autoDetectBuildpacks: boolean;
}> = ({
  porterApp,
  updatePorterApp,
  autoDetectBuildpacks,
}) => {
    const { currentProject } = useContext(Context);

    const [selectedStack, setSelectedStack] = useState<string>("");
    const [stackOptions, setStackOptions] = useState<{ label: string; value: string }[]>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isDetectingBuildpacks, setIsDetectingBuildpacks] = useState(false);
    const [error, setError] = useState<string>("");

    const [selectedBuildpacks, setSelectedBuildpacks] = useState<Buildpack[]>([]);
    const [availableBuildpacks, setAvailableBuildpacks] = useState<Buildpack[]>([]);

    const detectAndSetBuildPacks = async (detect: boolean) => {
      try {
        if (currentProject == null) {
          return;
        }

        if (!detect) {
          // in this case, we are not detecting buildpacks, so we just populate based on the DB
          if (porterApp.builder != null) {
            setSelectedStack(porterApp.builder);
            setStackOptions([{ label: porterApp.builder, value: porterApp.builder }]);
          }
          if (porterApp.buildpacks != null) {
            setSelectedBuildpacks(porterApp.buildpacks.map(bp => ({
              name: BUILDPACK_TO_NAME[bp] ?? bp,
              buildpack: bp,
              config: {},
            })));
          }
        } else {
          if (isDetectingBuildpacks) {
            return;
          }
          setIsDetectingBuildpacks(true);
          const detectBuildPackRes = await api.detectBuildpack(
            "<token>",
            {
              dir: porterApp.build_context || ".",
            },
            {
              project_id: currentProject.id,
              git_repo_id: porterApp.git_repo_id,
              kind: "github",
              owner: porterApp.repo_name.split("/")[0],
              name: porterApp.repo_name.split("/")[1],
              branch: porterApp.git_branch,
            }
          );

          const builders = detectBuildPackRes.data as DetectedBuildpack[];
          if (builders.length === 0) {
            return;
          }
          setStackOptions(builders.flatMap((builder) => {
            return builder.builders.map((stack) => ({
              label: `${builder.name} - ${stack}`,
              value: stack.toLowerCase(),
            }));
          }).sort((a, b) => {
            if (a.label < b.label) {
              return -1;
            }
            if (a.label > b.label) {
              return 1;
            }
            return 0;
          }));

          const defaultBuilder = builders.find(
            (builder) => builder.name.toLowerCase() === DEFAULT_BUILDER_NAME
          ) ?? builders[0];

          const allBuildpacks = defaultBuilder.others.concat(defaultBuilder.detected);

          let detectedBuilder: string;
          if (defaultBuilder.builders.length && defaultBuilder.builders.includes(DEFAULT_HEROKU_STACK)) {
            setSelectedStack(DEFAULT_HEROKU_STACK);
            detectedBuilder = DEFAULT_HEROKU_STACK;
          } else {
            setSelectedStack(defaultBuilder.builders[0]);
            detectedBuilder = defaultBuilder.builders[0];
          }

          const newBuildpacks = defaultBuilder.detected.filter(bp => !porterApp.buildpacks.includes(bp.buildpack));
          if (autoDetectBuildpacks) {
            updatePorterApp({ builder: detectedBuilder, buildpacks: [...porterApp.buildpacks, ...newBuildpacks.map(bp => bp.buildpack)] });
            setSelectedBuildpacks(defaultBuilder.detected);
            setAvailableBuildpacks(defaultBuilder.others);
            setError("");
          } else {
            updatePorterApp({ builder: detectedBuilder });
            setAvailableBuildpacks(allBuildpacks.filter(bp => !porterApp.buildpacks?.includes(bp.buildpack)));
          }
        }
      } catch (err) {
        setError(`Unable to detect buildpacks at path: ${porterApp.build_context}. Please make sure your repo, branch, and application root path are all set correctly and attempt to detect again.`);
      } finally {
        setIsDetectingBuildpacks(false);
      }
    }

    useEffect(() => {
      detectAndSetBuildPacks(autoDetectBuildpacks);
    }, [currentProject]);

    const handleAddCustomBuildpack = (buildpack: Buildpack) => {
      if (porterApp.buildpacks.find((bp) => bp === buildpack.buildpack) == null) {
        updatePorterApp({ buildpacks: [...porterApp.buildpacks, buildpack.buildpack] });
        setSelectedBuildpacks([...selectedBuildpacks, buildpack]);
      }
    };

    return (
      <BuildpackConfigurationContainer>
        {selectedBuildpacks.length > 0 && (
          <>
            <Helper>
              The following buildpacks were automatically detected. You can also
              manually add, remove, or re-order buildpacks here.
            </Helper>
            <BuildpackList
              selectedBuildpacks={selectedBuildpacks}
              setSelectedBuildpacks={setSelectedBuildpacks}
              availableBuildpacks={availableBuildpacks}
              setAvailableBuildpacks={setAvailableBuildpacks}
              porterApp={porterApp}
              updatePorterApp={updatePorterApp}
              showAvailableBuildpacks={false}
              isDetectingBuildpacks={isDetectingBuildpacks}
              detectBuildpacksError={error}
              droppableId={"non-modal"}
            />
          </>
        )}
        {autoDetectBuildpacks && error !== "" && (
          <>
            <Spacer y={1} />
            <Error message={error} />
          </>
        )}
        <Spacer y={1} />
        <Button onClick={() => {
          setIsModalOpen(true);
          setError("");
        }}>
          <I className="material-icons">add</I> Add / detect buildpacks
        </Button>
        {isModalOpen && (
          <BuildpackConfigurationModal
            closeModal={() => setIsModalOpen(false)}
            selectedStack={selectedStack}
            sortedStackOptions={stackOptions}
            setStackValue={(option) => {
              setSelectedStack(option);
              updatePorterApp({ builder: option });
            }}
            selectedBuildpacks={selectedBuildpacks}
            setSelectedBuildpacks={setSelectedBuildpacks}
            availableBuildpacks={availableBuildpacks}
            setAvailableBuildpacks={setAvailableBuildpacks}
            porterApp={porterApp}
            updatePorterApp={updatePorterApp}
            isDetectingBuildpacks={isDetectingBuildpacks}
            detectBuildpacksError={error}
            handleAddCustomBuildpack={handleAddCustomBuildpack}
            detectAndSetBuildPacks={detectAndSetBuildPacks}
          />
        )}
      </BuildpackConfigurationContainer>
    );
  };

export default BuildpackSettings;


const I = styled.i`
  color: white;
  font-size: 14px;
  display: flex;
  align-items: center;
  margin-right: 5px;
  justify-content: center;
`;

const fadeIn = keyframes`
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
`;

const BuildpackConfigurationContainer = styled.div`
  animation: ${fadeIn} 0.75s;
`;
