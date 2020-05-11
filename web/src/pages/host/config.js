
import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal,
 Select, Row, Col, Button, Popconfirm, Icon, 
 message} from "antd";
import {hasPermission} from "@/utils/globalTools"
import {connect} from "dva";

const FormItem = Form.Item;

@connect(({ loading, host, config}) => {
    return {
      hostAppList: host.hostAppList,
      hostAppLen: host.hostAppLen,
      hostAppLoading: loading.effects['host/getHostApp'],
      projectsList: config.projectsList,
    }
})

class HostAppPage extends React.Component {
  state = {
    visible: false,
    hid : 0,
  };

  componentDidMount() {
    const { location, dispatch } = this.props;
    this.setState({ 
      hid: location.search.split("=")[1] 
    });

    dispatch({
      type: 'config/getProject',
      payload: {
        page: 1,
        pageSize: 999,
        active: 1,
      }
    });

    dispatch({
      type: 'host/getHostApp',
      payload: {
        page: 1,
        pageSize: 10,
        hid: location.search.split("=")[1],
      }
    });
  }

  // 翻页
  pageChange = (page) => {
    const { dispatch, location } = this.props;
    dispatch({
      type: 'host/getHostApp',
      payload: {
        page: page,
        pageSize: 10,
        hid: location.search.split("=")[1],
      }
    });
  };

  showTypeAddModal = () => {
    this.setState({ 
      visible: true 
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
    });
  };

  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        this.props.projectsList.map(x => {
          if (x.Name === values.Aid ) {
            values.Aid = x.id
          }
        })

        values.Hid = parseInt(this.state.hid)
        dispatch({
          type: 'host/hostAppAdd',
          payload: values,
        });

        // 重置 `visible` 属性为 false 以关闭对话框
        this.setState({ visible: false });
      }
    });
  };

  cancel  = () => { 
  }


  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'host/hostAppDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  columns = [
  {
    title: '已绑定业务',
    dataIndex: 'Aid',
    'render': Aid => this.props.projectsList.map(x => Aid == x.id && x.Name)
  }, {
    title: '备注',
    dataIndex: 'Desc',
    ellipsis: true
  }, {
    title: '操作',
    width: 200,
    render: (text, record) => (
      <span>
          <Popconfirm title="你确定要删除吗?" 
            onConfirm={()=>{this.deleteRecord(record)}} 
            onCancel={this.cancel()}
          >
            {hasPermission('host-app-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
          </Popconfirm>
      </span>
    ),
  }];
  
  render() {
    const {visible, hid} = this.state;
    const {projectsList, hostAppList, hostAppLen, hostAppLoading, form: { getFieldDecorator } } = this.props;
    const addHostApp = <Button type="primary" onClick={this.showTypeAddModal} >绑定新业务</Button>;

    const extra = <Row gutter={16}>
          {hasPermission('host-app-add') && <Col span={10}>{addHostApp}</Col>}
      </Row>;
    return (
      <div>
        <Modal
          title= { "绑定新业务" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Form>
            <FormItem label="业务名字">
              {getFieldDecorator('Aid', {
                rules: [{ required: true }],
              })(
                <Select
                  showSearch
                  searchPlaceholder="输入关键词"
                  notFoundContent="无法找到"
                  placeholder="Please select"
                  style={{ width: '100%' }}
                  tokenSeparators={[',']}
                >
                  {projectsList.map(
                    x => <Select.Option key={x.id} value={x.Name}>{x.Name}</Select.Option>
                  )}
                </Select>
              )}
            </FormItem>
            <FormItem label="备注信息">
              {getFieldDecorator('Desc', {
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
          </Form> 
        </Modal>

        <Card title="" extra={extra}>
          <Table  
          pagination={{
            showQuickJumper: true,
            total: hostAppLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={hostAppList} loading={hostAppLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(HostAppPage);
