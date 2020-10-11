"""
Backfill state_id and make it unique.

Revision ID: d42c175b439a
Revises: 3e8cc74a1e7b
Create Date: 2017-01-18 15:11:01.635632
"""

# revision identifiers, used by Alembic.
revision = "d42c175b439a"
down_revision = "3e8cc74a1e7b"

import sqlalchemy as sa
from sqlalchemy.dialects import mysql


def upgrade(op, tables, tester):
    # Backfill the queueitem table's state_id field with unique values for all entries which are
    # empty.
    conn = op.get_bind()
    conn.execute("update queueitem set state_id = id where state_id = ''")

    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_index("queueitem_state_id", table_name="queueitem")
    op.create_index("queueitem_state_id", "queueitem", ["state_id"], unique=True)
    # ### end Alembic commands ###


def downgrade(op, tables, tester):
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_index("queueitem_state_id", table_name="queueitem")
    op.create_index("queueitem_state_id", "queueitem", ["state_id"], unique=False)
    # ### end Alembic commands ###