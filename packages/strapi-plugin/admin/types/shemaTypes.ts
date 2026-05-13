export interface SchemaObject {
  $schema: string;
  $id: string;
  title: string;
  type: string;
  required: string[];
  additionalProperties: boolean;
  properties: Properties;
  allOf: AllOfItem[];
}

export interface $schema {
  type: string;
  description: string;
}

export interface Name {
  type: string;
  description: string;
  pattern: string;
}

export interface Label {
  type: string;
  description: string;
  minLength: number;
}

export interface Kind {
  type: string;
  description: string;
  enum: string[];
}

export interface SchemaVersion {
  const: string;
  default: string;
  description: string;
}

export interface Category {
  type: string;
  description: string;
}

export interface Columns {
  type: string;
  minimum: number;
  maximum: number;
  description: string;
}

export interface FullWidth {
  type: string;
  description: string;
}

export interface Gap {
  type: string;
  enum: string[];
  description: string;
}

export interface Breakpoint {
  type: string;
  enum: string[];
  description: string;
}

export interface Behavior {
  type: string;
  enum: string[];
  description: string;
}

export interface Properties {
  $schema: $schema;
  name: Name;
  label: Label;
  kind: Kind;
  schemaVersion: SchemaVersion;
  category: Category;
  layoutConfig: LayoutConfig;
  blockConfig: BlockConfig;
  sectionConfig: SectionConfig;
  dataset: Dataset;
  renderMeta: RenderMeta;
}

export interface Items {
  type: string;
  required: string[];
  additionalProperties: boolean;
  properties: Properties;
}

export interface Responsive {
  type: string;
  description: string;
  items: Items;
}

export interface Zones {
  type: string;
  description: string;
  minItems: number;
  items: Items;
}

export interface LayoutConfig {
  type: string;
  description: string;
  required: string[];
  additionalProperties: boolean;
  properties: Properties;
}

export interface BlockConfig {
  type: string;
  description: string;
  required: string[];
  additionalProperties: boolean;
  properties: Properties;
}

export interface SectionConfig {
  type: string;
  description: string;
  required: string[];
  additionalProperties: boolean;
  properties: Properties;
}

export interface Dataset {
  type: string;
  description: string;
  minItems: number;
  items: Items;
}

export interface RenderMeta {
  type: string;
  required: string[];
  additionalProperties: boolean;
  properties: Properties;
}

export interface If {
  properties: Properties;
  required: string[];
}

export interface Then {
  required: string[];
}

export interface AllOfItem {
  if: If;
  then: Then;
}
