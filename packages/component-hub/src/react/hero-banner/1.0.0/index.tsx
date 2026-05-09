import type { JSX } from "react";

type HeroBannerProps = {
  headline: string;
  subheadline?: string;
};

export function HeroBanner({
  headline,
  subheadline,
}: HeroBannerProps): JSX.Element {
  return (
    <section>
      <h1>{headline}</h1>
      {subheadline ? <p>{subheadline}</p> : null}
    </section>
  );
}
