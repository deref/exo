export interface Route<TProps> {
  pattern: string;
  load: () => Promise<TProps>;
  render: (props: TProps) => React.ReactNode;
}

export const makeRoute = <TProps>(route: Route<TProps>): Route<TProps> => route;
