import React, { useEffect } from "react";

export default function Home() {
  useEffect(() => {
    document.location.replace("/guide");
  });

  return <div />;
}
