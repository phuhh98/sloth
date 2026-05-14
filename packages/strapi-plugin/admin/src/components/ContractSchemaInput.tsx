import { forwardRef, useEffect, useState } from 'react';
import { JSONInput, Field, Button } from '@strapi/design-system';

import { useIntl } from 'react-intl';
import styled from 'styled-components';

const JSONInputWrapper = styled.div`
  position: relative;
`;

const ValidateSchemaButton = styled(Button)`
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  z-index: 1;
`;

const ContractSchemaInput = forwardRef<HTMLInputElement, any>((props, ref) => {
  const { attribute, disabled, intlLabel, name, onChange, required } = props; // these are just some of the props passed by the content-manager

  //   const { formatMessage } = useIntl();
  useEffect(() => {
    console.log('ContractSchemaInput props:', props);
  }, []);

  const DEFAULT_EMPTY_VALUE = '\n\n';

  const [internalValue, setInternalValue] = useState(props.value ?? DEFAULT_EMPTY_VALUE);

  const handleChange: React.ComponentProps<typeof JSONInput>['onBlur'] = (value) => {
    setInternalValue(JSON.stringify(value, null, 3)); // format the JSON string with indentation
  };

  return (
    <>
      <Field.Root id="contract-schema-input" error={''} hint="Description line lorem ipsum">
        <Field.Label>{JSON.stringify(intlLabel)}</Field.Label>
        <JSONInputWrapper>
          <JSONInput
            ref={ref}
            disabled={disabled}
            value={internalValue}
            required={required}
            aria-label="JSON"
            onChange={(v) => setInternalValue(v)}
            onBlurCapture={handleChange}
          ></JSONInput>
          <ValidateSchemaButton onClick={() => alert('validate schema')}>
            Validate Schema
          </ValidateSchemaButton>
        </JSONInputWrapper>

        <Field.Error></Field.Error>
        <Field.Hint></Field.Hint>
      </Field.Root>
      {/* {formatMessage(intlLabel)} */}
      {/* {intlLabel?.defaultMessage} */}
    </>
  );
});

export default ContractSchemaInput;
