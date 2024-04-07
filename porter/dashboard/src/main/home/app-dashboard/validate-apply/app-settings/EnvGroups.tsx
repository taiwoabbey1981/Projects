import React, { useMemo, useState } from "react";
import { useFieldArray, useFormContext } from "react-hook-form";
import styled from "styled-components";
import { type IterableElement } from "type-fest";

import Icon from "components/porter/Icon";
import Spacer from "components/porter/Spacer";
import Text from "components/porter/Text";
import { type PorterAppFormData } from "lib/porter-apps";

import { valueExists } from "shared/util";
import doppler from "assets/doppler.png";
import sliders from "assets/sliders.svg";

import EnvGroupModal from "./EnvGroupModal";
import ExpandableEnvGroup from "./ExpandableEnvGroup";
import { type PopulatedEnvGroup } from "./types";

type Props = {
  baseEnvGroups?: PopulatedEnvGroup[];
  attachedEnvGroups?: PopulatedEnvGroup[];
};

const EnvGroups: React.FC<Props> = ({
  baseEnvGroups = [],
  attachedEnvGroups = [],
}) => {
  const [showEnvModal, setShowEnvModal] = useState(false);

  const { control } = useFormContext<PorterAppFormData>();
  const {
    append,
    remove,
    fields: envGroups,
  } = useFieldArray({
    control,
    name: "app.envGroups",
  });
  const {
    append: appendDeletion,
    remove: removeDeletion,
    fields: deletedEnvGroups,
  } = useFieldArray({
    control,
    name: "deletions.envGroupNames",
  });

  const populatedEnvWithFallback = useMemo(() => {
    return envGroups
      .map((envGroup, index) => {
        const attachedEnvGroup = attachedEnvGroups.find(
          (attachedEnvGroup) => attachedEnvGroup.name === envGroup.name
        );

        if (attachedEnvGroup) {
          return {
            id: envGroup.id,
            envGroup: attachedEnvGroup,
            index,
          };
        }

        const baseEnvGroup = baseEnvGroups.find(
          (baseEnvGroup) => baseEnvGroup.name === envGroup.name
        );

        if (baseEnvGroup) {
          return {
            id: envGroup.id,
            envGroup: baseEnvGroup,
            index,
          };
        }

        return undefined;
      })
      .filter(valueExists);
  }, [envGroups, attachedEnvGroups, baseEnvGroups]);

  const onAdd = (
    inp: IterableElement<PorterAppFormData["app"]["envGroups"]>
  ): void => {
    const previouslyDeleted = deletedEnvGroups.findIndex(
      (s) => s.name === inp.name
    );

    if (previouslyDeleted !== -1) {
      removeDeletion(previouslyDeleted);
    }

    append(inp);
  };

  const onRemove = (index: number): void => {
    const name = populatedEnvWithFallback[index].envGroup.name;
    remove(index);

    const existingEnvGroupNames = envGroups.map((eg) => eg.name);
    if (existingEnvGroupNames.includes(name)) {
      appendDeletion({ name });
    }
  };

  return (
    <div>
      <LoadButton
        disabled={false}
        onClick={() => {
          setShowEnvModal(true);
        }}
      >
        <img src={sliders} /> Load from Env Group
      </LoadButton>
      {populatedEnvWithFallback.length > 0 && (
        <>
          <Spacer y={0.5} />
          <Text size={16}>Synced environment groups</Text>
          {populatedEnvWithFallback.map(({ envGroup, id, index }) => {
            return (
              <ExpandableEnvGroup
                key={id}
                index={index}
                envGroup={envGroup}
                remove={onRemove}
                icon={
                  <Icon src={envGroup.type === "doppler" ? doppler : sliders} />
                }
              />
            );
          })}
        </>
      )}
      {showEnvModal ? (
        <EnvGroupModal
          setOpen={setShowEnvModal}
          baseEnvGroups={baseEnvGroups}
          append={onAdd}
        />
      ) : null}
    </div>
  );
};

export default EnvGroups;

const AddRowButton = styled.div`
  display: flex;
  align-items: center;
  width: 270px;
  font-size: 13px;
  color: #aaaabb;
  height: 32px;
  border-radius: 3px;
  cursor: pointer;
  background: #ffffff11;
  :hover {
    background: #ffffff22;
  }

  > i {
    color: #ffffff44;
    font-size: 16px;
    margin-left: 8px;
    margin-right: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
`;

const LoadButton = styled(AddRowButton)<{ disabled?: boolean }>`
  background: ${(props) => (props.disabled ? "#aaaaaa55" : "none")};
  border: 1px solid ${(props) => (props.disabled ? "#aaaaaa55" : "#ffffff55")};
  cursor: ${(props) => (props.disabled ? "not-allowed" : "pointer")};

  > i {
    color: ${(props) => (props.disabled ? "#aaaaaa44" : "#ffffff44")};
    font-size: 16px;
    margin-left: 8px;
    margin-right: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  > img {
    width: 14px;
    margin-left: 10px;
    margin-right: 12px;
    opacity: ${(props) => (props.disabled ? "0.5" : "1")};
  }
`;
