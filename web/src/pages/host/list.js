import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal,
 Select, Row, Col, Button, Popconfirm, Icon, InputNumber,
 Switch, message} from "antd";
import SearchForm from '@/components/SearchForm';
import Import from './Import';
import {connect} from "dva";
import Link from 'umi/link';
import {hasPermission} from "@/utils/globalTools"

const FormItem = Form.Item;

const hostsZone = [
  {id: 1, name: '阿里云-华北2-C'},
  {id: 2, name: '阿里云-华北2-E'},
  {id: 3, name: 'IDC机房1'},
];

@connect(({ loading, host, config }) => {
    return {
      hostRoleList: host.hostRoleList,
      hostList: host.hostList,
      hostLen: host.hostLen,
      hostLoading: loading.effects['host/getHost'],
      configEnvList: config.configEnvList,
    }
 })


class HostPage extends React.Component {
  state = {
    visible: false,
    importVisible: false,
    editCacheData: {},
    Rid: 0,
    Name: '',
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ type: 'host/getHost',
      payload: {
        page: 1,
        pageSize: 10, 
      }
    });
    dispatch({ 
      type: 'host/getHostRole',
      payload: {
        page: 1,
        pageSize: 100, 
      }
    });
    dispatch({
      type: 'config/getConfigEnv',
      payload: {
        page: 1,
        pageSize: 50, 
      }
    });
  }
  
  showHostAddModal = () => {
    this.setState({ 
      editCacheData: {},
      visible: true 
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
      importVisible: false,
    });
  };

  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        values.Enable = values.Enable && 1 || 0
        if (Object.keys(obj).length) {
          if (
            obj.Rid      === values.Rid && 
            obj.ZoneId   === values.ZoneId && 
            obj.Enable   === values.Enable && 
            obj.Name     === values.Name && 
            obj.Username === values.Username && 
            obj.Addres   === values.Addres && 
            obj.EnvId    === values.EnvId && 
            obj.Port     === values.Port && 
            obj.Desc     === values.Desc 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id
            values.Port = parseInt(values.Port)
            dispatch({
              type: 'host/hostEdit',
              payload: values,
            });
          }
        } else {
          values.Status = 1
          values.Port = parseInt(values.Port)
          dispatch({
            type: 'host/hostAdd',
            payload: values,
          });
        }
        // 重置 `visible` 属性为 false 以关闭对话框
        this.setState({ visible: false });
      }
    });
  };

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'host/hostDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  // Popconfirm 取消事件
  cancel = () => {
  };

  // 禁用/启用主机
  changeActive = (values) => {
    const { dispatch } = this.props;
    if (values.Status) {
      values.Status = 0
    } else {
      values.Status = 1
    }

    dispatch({
      type: 'host/hostEdit',
      payload: values,
    })
  };

  //显示编辑界面
  handleEdit = (values, page) => {
    values.title =  '编辑主机-' + values.Name;
    this.setState({ 
      visible: true ,
      editCacheData: values
    });
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'host/getHost',
      payload: {
        page: page,
        pageSize: 10, 
      }
    });
  };

  fetchRecords = () => {
    const { dispatch } = this.props;
    dispatch({ type: 'host/getHost',
      payload: {
        page: 1,
        pageSize: 10, 
        Rid : this.state.Rid,
        Name: this.state.Name,
      }
    });
  };

  handleConsole = (info) => {
    window.open(`/admin/host/ssh/${info.id}?x-token=${sessionStorage.getItem('jwt')}`)
  };


  columns = [{
    title: '序号',
    key: 'series',
    render: (_, __, index) => index + 1,
    width: 80
  }, {
    title: '主机类别',
    dataIndex: 'Rid',
    'render': Rid => this.props.hostRoleList.map(x => Rid == x.id && x.Name)
  }, {
    title: '区域',
    dataIndex: 'ZoneId',
    'render': ZoneId => hostsZone.map(x => {
      if (ZoneId == x.id) {
        return x.name
      }
    })
  }, {
    title: '环境',
    dataIndex: 'EnvId',
    'render': EnvId =>  this.props.configEnvList.map(x => {
      if (EnvId==x.id) {
        return x.Name
      }
    })
  }, {
    title: '主机别名',
    dataIndex: 'Name',
  }, {
    title: '主机地址',
    dataIndex: 'Addres',
  },  {
    title: '端口',
    dataIndex: 'Port',
  },{
    title: '备注',
    dataIndex: 'Desc',
    ellipsis: true
  }, {
    title: '操作',
    width: 200,
    render: (text, record) => (
      <span>
        {
          hasPermission('host-edit') && 
          <a onClick={()=>{this.handleEdit(record)}}>
            <Icon type="edit"/>编辑
          </a>
        }
        <Divider type="vertical" />
        {  
          hasPermission('host-config') && 
          <Link to={`/host/config?hid=${record.id}`}>
            <Icon type="redo"/>业务分配
          </Link>
        }
        <Divider type="vertical" />
        {
          record.Status
          && 
          <Popconfirm title="你确定要禁用主机吗?" 
            onConfirm={()=>{this.changeActive(record)}} 
            onCancel={()=>{this.cancel()}}>
            {hasPermission('host-edit') && <a title="lock" ><Icon type="lock"/>禁用主机</a>}
          </Popconfirm> 
          || 
          <Popconfirm title="你确定要启主机吗?" 
            onConfirm={()=>{this.changeActive(record)}} 
            onCancel={()=>{this.cancel()}}>
            {hasPermission('host-edit') && <a title="删除" ><Icon type="unlock"/>启用主机</a>}
          </Popconfirm> 
        }
        <Divider type="vertical" />
          <Popconfirm title="你确定要删除吗?" 
            onConfirm={()=>{this.deleteRecord(record.id)}} 
            onCancel={()=>{this.cancel()}}>
            {hasPermission('host-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
          </Popconfirm>
        <Divider type="vertical" />
          <a auth="host.host.console" onClick={() => this.handleConsole(record)}>Console</a>
      </span>
    ),
  }];
  
  render() {
    const {visible, importVisible, editCacheData, Rid, Name} = this.state;
    const {configEnvList, hostRoleList, hostList, hostLen, 
      hostLoading, form: { getFieldDecorator } } = this.props;
    const addHost = <Button type="primary" onClick={this.showHostAddModal} >新增主机</Button>;
    const improtHosts = <Button style={{marginLeft: 20}} type="primary" icon="import"
                onClick={() => this.setState({importVisible: true})}>批量导入</Button>
    const extra = <Row gutter={16}>
        {hasPermission('host-add') && <Col span={10}>{addHost}</Col>}
        {hasPermission('host-import') && <Col span={10}>{improtHosts}</Col>}
      </Row>;

    return (
      <div>
        <SearchForm>
          <SearchForm.Item span={8} title="主机类别">
                <Select
                  placeholder="请选择"
                  onChange={value => this.state.Rid = value}
                  style={{ width: '100%' }}
                >
                {hostRoleList.map(x => 
                  <Select.Option key={x.id} value={x.id}>
                    {x.Name}
                  </Select.Option>)}
              </Select>
          </SearchForm.Item>
          <SearchForm.Item span={8} title="主机别名">
            <Input allowClear  onChange={e => this.state.Name = e.target.value} placeholder="请输入"/>
          </SearchForm.Item>
          <SearchForm.Item span={8}>
            <Button type="primary" icon="sync" onClick={this.fetchRecords}>搜索</Button>
          </SearchForm.Item>
        </SearchForm>
        
        <Modal
          width={800}
          maskClosable={false}
          title= {editCacheData.title || "新增主机" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
            <FormItem label="主机类型">
              {getFieldDecorator('Rid', {
                initialValue: editCacheData.Rid || 'Please select' ,
                rules: [{ required: true }],
              })(
                <Select
                  placeholder="Please select"
                  style={{ width: '100%' }}
                >
                  {hostRoleList.map(
                    x => <Select.Option key={x.id} value={x.id}>{x.Name}</Select.Option>
                  )}
                </Select>
              )}
            </FormItem>
            <FormItem label="主机区域">
               {getFieldDecorator('ZoneId', {
                initialValue: editCacheData.ZoneId || 'Please select' ,
                rules: [{ required: true }],
              })( 
                <Select
                  placeholder="Please select"
                >
                {hostsZone.map(x => <Select.Option key={x.id} value={x.id}>{x.name}</Select.Option>)}
                </Select>
              )}
            </FormItem>
            <FormItem label="主机环境">
               {getFieldDecorator('EnvId', {
                initialValue: editCacheData.EnvId || 'Please select' ,
                rules: [{ required: true }],
              })( 
                <Select
                  placeholder="Please select"
                  style={{ width: '100%' }}
                >
                  {configEnvList.map(x => 
                  <Select.Option key={x.id} value={x.id}>{x.Name}</Select.Option>)}
                </Select>
              )}
            </FormItem>
            <FormItem label="主机别名">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <Form.Item required label="连接地址" style={{marginBottom: 0}}>
              <FormItem style={{display: 'inline-block', width: 'calc(30%)'}}>
                {getFieldDecorator('Username', {
                  initialValue: editCacheData.Username || '',
                  rules: [{ required: true }],
                })(
                  <Input addonBefore="ssh" placeholder="用户名"/>
                )}
              </FormItem>
              <FormItem style={{display: 'inline-block', width: 'calc(40%)'}}>
                {getFieldDecorator('Addres', {
                  initialValue: editCacheData.Addres || '',
                  rules: [{ required: true }],
                })(
                  <Input addonBefore="@" placeholder="主机名/IP"/>
                )}
              </FormItem>
              <FormItem style={{display: 'inline-block', width: 'calc(30%)'}}>
                {getFieldDecorator('Port', {
                  initialValue: editCacheData.Port || '',
                  rules: [{ required: true }],
                })(
                  <Input addonBefore="-p" placeholder="端口"/>
                )}
              </FormItem>
            </Form.Item>
            <FormItem label="跳板机接管">
              {getFieldDecorator('Enable', {
                initialValue: editCacheData.Enable || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.Enable || false} onChange={this.onCheckChange} />,
              )}
            </FormItem>
            <FormItem label="备注信息">
              {getFieldDecorator('Desc', {
                initialValue: editCacheData.Desc || '',
                rules: [{ required: true }],
              })(
                <Input.TextArea />
              )}
            </FormItem>
            { !Object.keys(editCacheData).length  && 
              <FormItem label="密码">
                {getFieldDecorator('Password', {
                  initialValue: editCacheData.Password || '',
                  rules: [{ required: true }],
                })(
                  <Input.Password />
                )}
              </FormItem> 
            }
            { !Object.keys(editCacheData).length  && 
              <FormItem wrapperCol={{span: 14, offset: 6}}>
                <span role="img" aria-label="notice">⚠️ 首次验证时需要登录用户名对应的密码，但不会存储该密码。</span>
              </FormItem>
            }
          </Form> 
        </Modal>

        { importVisible && 
          <Import 
            onCancel={this.handleCancel}
            dispatch={this.props.dispatch}
          />
        }

        <Card title="" extra={extra}>
          <Table  
          pagination={{
            defaultPageSize: 10,
            showQuickJumper: true,
            total: hostLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={hostList} loading={hostLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(HostPage);

