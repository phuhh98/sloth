import { schemas } from '@sloth/contracts';
import { Ajv } from 'ajv';

const COMPONENT_CONTRACT_SCHEMA = schemas.componentContract;

const ajv = new Ajv();
const componentContractValidator = ajv.compile(COMPONENT_CONTRACT_SCHEMA);

export function validateComponentContract(contract: unknown) {
  const isValid = componentContractValidator(contract);
  if (!isValid) {
    return {
      valid: false,
      errors: componentContractValidator.errors,
    };
  }

  return {
    valid: true,
    errors: null,
  };
}
