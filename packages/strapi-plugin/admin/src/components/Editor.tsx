import { Puck, type Config } from '@puckeditor/core';
import puckCss from '@puckeditor/core/dist/index.css?inline';
import { useEffect } from 'react';

// Create Puck component config
const config: Config = {
  categories: {
    layout: { components: ['HeadingBlock', 'MainAside'], title: 'Layout' },
    block: {
      components: ['CTAButton'],
      title: 'Block',
    },
  },
  components: {
    HeadingBlock: {
      fields: {
        children: {
          type: 'text',
        },
      },
      render: ({ children }) => {
        return <h1>{children}</h1>;
      },
    },

    MainAside: {
      fields: {
        content: {
          type: 'slot',
        },
        aside: {
          type: 'slot',
        },
        backgroundColor: {
          type: 'select',
          options: [
            { value: 'white', label: 'White' },
            { value: 'lightgray', label: 'Light Gray' },
            { value: 'lightblue', label: 'Light Blue' },
          ],
        },
      },

      render: ({ content: Content, aside: Aside, backgroundColor }) => {
        return (
          <div
            style={{
              display: 'grid',
              gap: '2rem',
              gridTemplateColumns: '1fr '.repeat(12),
              backgroundColor,
              margin: '0 auto',
            }}
          >
            <Content
              style={{
                gridColumn: 'span 8',
                display: 'grid',
                gap: '1rem',
                gridTemplateColumns: '1fr '.repeat(8),
              }}
            />
            <Aside
              style={{
                gridColumn: 'span 4',
                display: 'grid',
                gap: '1rem',
                gridTemplateColumns: '1fr '.repeat(4),
              }}
            />
          </div>
        );
      },
    },

    CTAButton: {
      inline: true,
      fields: {
        label: {
          type: 'text',
        },
        url: {
          type: 'text',
        },
      },

      render: ({ label, url, puck }) => {
        return (
          <div
            ref={puck.dragRef}
            style={{
              width: '100%',
              gridColumn: 'span 4',
            }}
          >
            <a href={url} target="_blank" rel="noopener noreferrer">
              <button>{label}</button>
            </a>
          </div>
        );
      },
    },
  },
};

// Describe the initial data
const initialData = {} as const;

// Save the data to your database
const save = (data: typeof initialData) => {};

// Render Puck editor
export function Editor() {
  useEffect(() => {
    const style = document.createElement('style');
    style.dataset['puckStyles'] = '';
    style.textContent = puckCss;
    document.head.appendChild(style);
    return () => {
      style.remove();
    };
  }, []);

  return <Puck config={config} data={initialData} onPublish={save} />;
}
