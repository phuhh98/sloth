import { Main } from '@strapi/design-system';
import { useIntl } from 'react-intl';

import { getTranslation } from '../utils/getTranslation';
import { Editor } from '../components/Editor';

const HomePage = () => {
  const { formatMessage } = useIntl();

  return (
    <Main>
      <Editor />
    </Main>
  );
};

export { HomePage };
