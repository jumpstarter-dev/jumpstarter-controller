#!/usr/bin/env -S uv run --script
# /// script
# requires-python = ">=3.12"
# dependencies = ["pydantic"]
# ///

from __future__ import annotations

import json

from enum import Enum
from typing import List, Optional, Union

from pydantic import BaseModel, ConfigDict, Field, RootModel, conint


class Provisioning(BaseModel):
    model_config = ConfigDict(extra="forbid")

    enabled: Optional[bool] = None


class Internal(BaseModel):
    model_config = ConfigDict(extra="forbid")

    prefix: Optional[str] = None


class Keepalive(BaseModel):
    model_config = ConfigDict(extra="forbid")

    # EnforcementPolicy parameters
    minTime: Optional[str] = Field(
        None,
        description="The minimum amount of time a client should wait before sending a keepalive ping",
    )
    permitWithoutStream: Optional[bool] = Field(
        None,
        description="Whether to allow keepalive pings even when there are no active streams(RPCs)",
    )

    # ServerParameters for connection timeout control
    timeout: Optional[str] = Field(
        None,
        description="How long the server waits for a ping response before closing the connection",
    )
    maxConnectionIdle: Optional[str] = Field(
        None,
        description="Maximum time a connection can be idle before being closed",
    )
    maxConnectionAge: Optional[str] = Field(
        None,
        description="Maximum lifetime of a connection before it's closed",
    )
    maxConnectionAgeGrace: Optional[str] = Field(
        None,
        description="Grace period after max connection age before forcible closure",
    )
    time: Optional[str] = Field(
        None,
        description="How often the server sends keepalive pings to clients",
    )


class Grpc(BaseModel):
    model_config = ConfigDict(extra="forbid")

    keepalive: Optional[Keepalive] = None


class Metrics(BaseModel):
    enabled: Optional[bool] = None


class Global(BaseModel):
    baseDomain: Optional[str] = Field(
        None, description="Base domain to construct the FQDN for the service endpoints"
    )
    metrics: Optional[Metrics] = None


class Mode(Enum):
    ingress = "ingress"
    route = "route"
    nodeport = "nodeport"


class Mode1(Enum):
    passthrough = "passthrough"
    reencrypt = "reencrypt"


class Port(RootModel):
    root: conint(ge=0, le=65535)


class Ingress(BaseModel):
    model_config = ConfigDict(extra="forbid")

    enabled: Optional[bool] = Field(
        None, description="Whether to enable Ingress for the gRPC endpoint"
    )
    class_: Optional[str] = Field(
        None, alias="class", description="IngressClass to use for the gRPC endpoint"
    )


class Route(BaseModel):
    model_config = ConfigDict(extra="forbid")

    enabled: Optional[bool] = Field(
        None, description="Whether to enable OpenShift Router for the gRPC endpoint"
    )


class PrefixedClaimOrExpression1(BaseModel):
    model_config = ConfigDict(extra="forbid")

    claim: str
    prefix: str


class PrefixedClaimOrExpression2(BaseModel):
    model_config = ConfigDict(extra="forbid")

    expression: str


class PrefixedClaimOrExpression(RootModel):
    root: Union[PrefixedClaimOrExpression1, PrefixedClaimOrExpression2]


class ClaimOrExpression1(BaseModel):
    model_config = ConfigDict(extra="forbid")

    claim: str
    expression: Optional[str] = None


class ClaimOrExpression2(BaseModel):
    model_config = ConfigDict(extra="forbid")

    claim: Optional[str] = None
    expression: str


class ClaimOrExpression(RootModel):
    root: Union[ClaimOrExpression1, ClaimOrExpression2]


class AudienceMatchPolicy(Enum):
    MatchAny = "MatchAny"


class Issuer(BaseModel):
    model_config = ConfigDict(extra="forbid")

    url: Optional[str] = None
    discoveryURL: Optional[str] = None
    certificateAuthority: Optional[str] = None
    audiences: Optional[List[str]] = None
    audienceMatchPolicy: Optional[AudienceMatchPolicy] = None


class ClaimValidationRule(BaseModel):
    model_config = ConfigDict(extra="forbid")

    claim: Optional[str] = None
    requiredValue: Optional[str] = None
    expression: Optional[str] = None
    message: Optional[str] = None


class ExtraItem(BaseModel):
    model_config = ConfigDict(extra="forbid")

    key: Optional[str] = None
    valueExpression: Optional[str] = None


class ClaimMappings(BaseModel):
    model_config = ConfigDict(extra="forbid")

    username: Optional[PrefixedClaimOrExpression] = None
    groups: Optional[PrefixedClaimOrExpression] = None
    uid: Optional[ClaimOrExpression] = None
    extra: Optional[List[ExtraItem]] = None


class UserValidationRule(BaseModel):
    model_config = ConfigDict(extra="forbid")

    expression: Optional[str] = None
    message: Optional[str] = None


class JWTAuthenticator(BaseModel):
    model_config = ConfigDict(extra="forbid")

    issuer: Optional[Issuer] = None
    claimValidationRules: Optional[List[ClaimValidationRule]] = None
    claimMappings: Optional[ClaimMappings] = None
    userValidationRules: Optional[List[UserValidationRule]] = None


class ExporterOptions(BaseModel):
    model_config = ConfigDict(extra="forbid")

    offlineTimeout: Optional[str] = Field(
        None,
        description="How long to wait before marking the exporter as offline",
    )


class Authentication(BaseModel):
    model_config = ConfigDict(extra="forbid")

    internal: Optional[Internal] = None
    jwt: Optional[List[JWTAuthenticator]] = Field(
        None,
        description="External OIDC authentication, see https://kubernetes.io/docs/reference/access-authn-authz/authentication/#using-authentication-configuration for documentation",
    )


class JumpstarterConfig(BaseModel):
    model_config = ConfigDict(extra="forbid")

    provisioning: Optional[Provisioning] = None
    authentication: Optional[Authentication] = None
    grpc: Optional[Grpc] = None
    exporterOptions: Optional[ExporterOptions] = None


class Nodeport(BaseModel):
    model_config = ConfigDict(extra="forbid")

    enabled: Optional[bool] = None
    port: Optional[Port] = None
    routerPort: Optional[Port] = None


class Tls(BaseModel):
    model_config = ConfigDict(extra="forbid")

    enabled: Optional[bool] = None
    secret: Optional[str] = None
    controllerCertSecret: Optional[str] = Field(
        None,
        description="Secret containing the TLS certificate/key for the gRPC controller endpoint",
    )
    routerCertSecret: Optional[str] = Field(
        None,
        description="Secret containing the TLS certificate/key for the gRPC router endpoints",
    )
    port: Optional[Port] = Field(
        None,
        description="Port to use for the gRPC endpoints Ingress or Route, this can be useful for ingress routers on non-standard ports",
    )
    mode: Optional[Mode1] = Field(None, description="TLS mode for gRPC endpoints")


class Grpc1(BaseModel):
    model_config = ConfigDict(extra="forbid")

    hostname: Optional[str] = Field(
        None, description="Hostname for the controller to use for the controller gRPC"
    )
    routerHostname: Optional[str] = Field(
        None, description="Hostname for the router to use for the router gRPC"
    )
    endpoint: Optional[str] = Field(
        None,
        description="The endpoints are passed down to the services to know where to announce the endpoints to the clients",
    )
    routerEndpoint: Optional[str] = Field(
        None,
        description="The endpoints are passed down to the services to know where to announce the endpoints to the clients",
    )
    additionalRouters: dict[str, Router] | None = Field(
        None, description="Additional routers to deploy"
    )
    ingress: Optional[Ingress] = None
    route: Optional[Route] = None
    nodeport: Optional[Nodeport] = None
    mode: Optional[Mode] = None
    tls: Optional[Tls] = None


class Router(BaseModel):
    model_config = ConfigDict(extra="forbid")
    hostname: str | None = None
    endpoint: str | None = None
    labels: dict[str, str] | None = None
    nodeSelector: dict[str, str] | None = None


class Model(BaseModel):
    model_config = ConfigDict(extra="forbid")

    enabled: Optional[bool] = Field(
        None, description="Whether to enable jumpstarter controller"
    )
    authenticationConfig: Optional[str] = None
    config: Optional[JumpstarterConfig] = None
    namespace: Optional[str] = Field(
        None,
        description="Namespace where the controller will be deployed, defaults to global.namespace",
    )
    image: str = Field(..., description="Image for the controller")
    tag: Optional[str] = Field(None, description="Image tag for the controller")
    imagePullPolicy: str = Field(
        ..., description="Image pull policy for the controller"
    )
    global_: Optional[Global] = Field(
        None, alias="global", description="Global parameters"
    )
    grpc: Optional[Grpc1] = None


print(json.dumps(Model.model_json_schema(), indent=2))
