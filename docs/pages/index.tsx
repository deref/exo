import React, { useEffect } from "react";

export default function Home() {
  useEffect(() => {
    document.location.replace("https://exo.deref.io");
  });

  return <div />;
}
