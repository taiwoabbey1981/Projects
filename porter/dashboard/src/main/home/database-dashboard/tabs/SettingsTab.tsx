import React, { useContext } from "react";
import styled from "styled-components";

import Button from "components/porter/Button";
import Icon from "components/porter/Icon";
import Spacer from "components/porter/Spacer";
import Text from "components/porter/Text";
import { useDatastoreMethods } from "lib/hooks/useDatabaseMethods";

import { Context } from "shared/Context";
import trash from "assets/trash.png";

import { useDatastoreContext } from "../DatabaseContextProvider";

const SettingsTab: React.FC = () => {
  const { setCurrentOverlay } = useContext(Context);
  const { datastore } = useDatastoreContext();
  const { deleteDatastore } = useDatastoreMethods();
  const handleDeletionSubmit = async (): Promise<void> => {
    if (setCurrentOverlay == null) {
      return;
    }

    try {
      await deleteDatastore(datastore.name);
      setCurrentOverlay(null);
    } catch (error) {
      // todo: handle error
    }
  };

  const handleDeletionClick = async (): Promise<void> => {
    if (setCurrentOverlay === undefined) {
      return;
    }

    setCurrentOverlay({
      message: `Are you sure you want to delete ${datastore.name}?`,
      onYes: handleDeletionSubmit,
      onNo: () => {
        setCurrentOverlay(null);
      },
    });
  };

  return (
    <StyledTemplateComponent>
      <InnerWrapper>
        <Text size={16}>Delete &quot;{datastore.name}&quot;</Text>
        <Spacer y={0.5} />
        <Text color="helper">
          Delete this datastore and all of its resources.
        </Text>
        <Spacer y={0.5} />
        <Button color="#b91133" onClick={handleDeletionClick}>
          <Icon src={trash} height={"15px"} />
          <Spacer inline x={0.5} />
          Delete {datastore.name}
        </Button>
      </InnerWrapper>
    </StyledTemplateComponent>
  );
};

export default SettingsTab;

const StyledTemplateComponent = styled.div`
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

const InnerWrapper = styled.div<{ full?: boolean }>`
  width: 100%;
  height: ${(props) => (props.full ? "100%" : "calc(100% - 65px)")};
  padding: 30px;
  padding-bottom: 15px;
  position: relative;
  overflow: auto;
  margin-bottom: 30px;
  border-radius: 5px;
  background: ${(props) => props.theme.fg};
  border: 1px solid #494b4f;
`;
