# -- 📖 Common import
from fastapi import FastAPI
from fastapi.security import OAuth2PasswordBearer
from sqlalchemy import create_engine
from helpers.transport import RMQEventMap, KafkaEventMap, create_kafka_connection_parameters, create_rmq_connection_parameters
from helpers.app import get_abs_static_dir, create_message_bus, create_rpc, handle_app_shutdown, register_static_dir_route_handler
from repos.dbUser import DBUserRepo
from auth import register_auth_route_handler, register_auth_event_handler, register_auth_rpc_handler, AuthModel, TokenModel, UserModel, UserSeederModel
from schemas.user import UserData

import os

# -- 🐇 Rabbitmq setting
rmq_connection_parameters = create_rmq_connection_parameters(
    host = os.getenv('APP_RABBITMQ_HOST', 'localhost'),
    user = os.getenv('APP_RABBITMQ_USER', 'root'),
    password = os.getenv('APP_RABBITMQ_PASS', ''),
    virtual_host = os.getenv('APP_RABBITMQ_VHOST', '/'),
    heartbeat=30
)
rmq_event_map = RMQEventMap({})

# -- 🪠 Kafka setting
kafka_connection_parameters = create_kafka_connection_parameters(
    bootstrap_servers = os.getenv('APP_KAFKA_BOOTSTRAP_SERVERS', 'localhost:29092'),
    sasl_mechanism = os.getenv('APP_KAFKA_SASL_MECHANISM', 'PLAIN'),
    sasl_plain_username = os.getenv('APP_KAFKA_SASL_PLAIN_USERNAME', 'root'),
    sasl_plain_password = os.getenv('APP_KAFKA_SASL_PLAIN_PASSWORD', '')
)
kafka_event_map = KafkaEventMap({})

# -- 🚌 Message bus and RPC initialization
mb_type = os.getenv('APP_MESSAGE_BUS_TYPE', 'local')
rpc_type = os.getenv('APP_RPC_TYPE', 'local')
mb = create_message_bus(mb_type, rmq_connection_parameters, rmq_event_map, kafka_connection_parameters, kafka_event_map)
rpc = create_rpc(rpc_type, rmq_connection_parameters, rmq_event_map)

# -- 🛢️ Database engine initialization
db_url = os.getenv('APP_SQLALCHEMY_DATABASE_URL', 'sqlite://')
engine = create_engine(db_url, echo=True)
user_repo = DBUserRepo(engine=engine, create_all=True)

# -- 👤 User initialization
guest_username = os.getenv('APP_GUEST_USERNAME', 'guest')
root_role = os.getenv('APP_ROOT_ROLE', 'root')
token_model = TokenModel(
    access_token_secret_key = os.getenv('APP_ACCESS_TOKEN_SECRET_KEY', '123'),
    access_token_algorithm = os.getenv('APP_ACCESS_TOKEN_ALGORITHM', 'HS256'),
    access_token_expire_minutes = int(os.getenv('APP_ACCESS_TOKEN_EXPIRE_MINUTES', '30'))
)
user_model = UserModel(user_repo, token_model, guest_username)
user_seeder_model = UserSeederModel(user_model)
user_seeder_model.seed(UserData(
    username = os.getenv('APP_ROOT_USERNAME', 'root'),
    email = os.getenv('APP_ROOT_INITIAL_EMAIL', 'root@root.com'),
    password = os.getenv('APP_ROOT_INITIAL_PASSWORD', 'toor'),
    active = True,
    roles = root_role,
    full_name = os.getenv('APP_ROOT_INITIAL_FULL_NAME', 'root')
))
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")
auth_model = AuthModel(user_model, oauth2_scheme, root_role)

# -- ⚡FastAPI initialization
app = FastAPI(title='ztplAppName')

# -- 📖 Register core handlers
enable_route_handler = os.getenv('APP_ENABLE_ROUTE_HANDLER', '1') != '0'
enable_event_handler = os.getenv('APP_ENABLE_EVENT_HANDLER', '1') != '0'
enable_rpc_handler = os.getenv('APP_ENABLE_RPC_HANDLER', '1') != '0'
static_url = os.getenv('APP_STATIC_URL', '/static')
static_dir = get_abs_static_dir(os.getenv('APP_STATIC_DIR', ''))
handle_app_shutdown(app, mb, rpc)
register_static_dir_route_handler(app, static_url, static_dir, static_route_name='static')
if enable_route_handler:
    register_auth_route_handler(app, mb, rpc, auth_model, user_model)
if enable_event_handler:
    register_auth_event_handler(mb)
if enable_rpc_handler:
    register_auth_rpc_handler(rpc, user_model)